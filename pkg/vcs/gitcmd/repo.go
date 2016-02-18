package gitcmd

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/binary"
	"github.com/golang/groupcache/lru"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sqs/pbtypes"
	"src.sourcegraph.com/sourcegraph/pkg/gitserver"
	"src.sourcegraph.com/sourcegraph/pkg/synclru"
	"src.sourcegraph.com/sourcegraph/pkg/vcs"
	"src.sourcegraph.com/sourcegraph/pkg/vcs/internal"
	"src.sourcegraph.com/sourcegraph/pkg/vcs/util"
)

var (
	// logEntryPattern is the regexp pattern that matches entries in the output of
	// the `git shortlog -sne` command.
	logEntryPattern = regexp.MustCompile(`^\s*([0-9]+)\s+([A-Za-z]+(?:\s[A-Za-z]+)*)\s+<([A-Za-z@.]+)>\s*$`)
)

type Repository struct {
	Dir string

	editLock   sync.RWMutex // protects ops that change repository data
	AppdashRec *appdash.Recorder
}

func (r *Repository) String() string {
	return fmt.Sprintf("git (cmd) repo at %s", r.Dir)
}

func Open(dir string) *Repository {
	return &Repository{Dir: dir}
}

// CloneOpt configures a clone operation.
type CloneOpt struct {
	Bare   bool // create a bare repo
	Mirror bool // create a mirror repo (`git clone --mirror`)

	vcs.RemoteOpts // configures communication with the remote repository
}

func Clone(url, dir string, opt CloneOpt) error {
	args := []string{"clone"}
	if opt.Bare {
		args = append(args, "--bare")
	}
	if opt.Mirror {
		args = append(args, "--mirror")
	}
	args = append(args, "--", url, filepath.ToSlash(dir))
	cmd := exec.Command("git", args...)

	if opt.SSH != nil {
		gitSSHWrapper, gitSSHWrapperDir, keyFile, err := makeGitSSHWrapper(opt.SSH.PrivateKey)
		defer func() {
			if keyFile != "" {
				if err := os.Remove(keyFile); err != nil {
					log.Fatalf("Error removing SSH key file %s: %s.", keyFile, err)
				}
			}
		}()
		if err != nil {
			return err
		}
		defer os.Remove(gitSSHWrapper)
		if gitSSHWrapperDir != "" {
			defer os.RemoveAll(gitSSHWrapperDir)
		}
		cmd.Env = []string{"GIT_SSH=" + gitSSHWrapper}
	}

	if opt.HTTPS != nil {
		env := environ(os.Environ())
		env.Unset("GIT_TERMINAL_PROMPT")

		gitPassHelper, gitPassHelperDir, err := makeGitPassHelper(opt.HTTPS.Pass)
		if err != nil {
			return err
		}
		defer os.Remove(gitPassHelper)
		if gitPassHelperDir != "" {
			defer os.RemoveAll(gitPassHelperDir)
		}
		env.Unset("GIT_ASKPASS")
		env = append(env, "GIT_ASKPASS="+gitPassHelper)

		cmd.Env = env
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec `git clone` failed: %s. Output was:\n\n%s", err, out)
	}
	return nil
}

// checkSpecArgSafety returns a non-nil err if spec begins with a "-", which could
// cause it to be interpreted as a git command line argument.
func checkSpecArgSafety(spec string) error {
	if strings.HasPrefix(spec, "-") {
		return errors.New("invalid git revision spec (begins with '-')")
	}
	return nil
}

// dividedOutput runs the command and returns its standard output and standard error.
func dividedOutput(c *exec.Cmd) (stdout []byte, stderr []byte, err error) {
	var outb, errb bytes.Buffer
	c.Stdout = &outb
	c.Stderr = &errb
	err = c.Run()
	return outb.Bytes(), errb.Bytes(), err
}

func (r *Repository) ResolveRevision(spec string) (vcs.CommitID, error) {
	defer r.trace(time.Now(), "ResolveRevision", spec)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	if err := checkSpecArgSafety(spec); err != nil {
		return "", err
	}

	cmd := gitserver.Command("git", "rev-parse", spec+"^0")
	cmd.Dir = r.Dir
	stdout, stderr, err := cmd.DividedOutput()
	if err != nil {
		if err == vcs.ErrRepoNotExist {
			return "", err
		}
		if bytes.Contains(stderr, []byte("unknown revision")) {
			return "", vcs.ErrRevisionNotFound
		}
		return "", fmt.Errorf("exec `git rev-parse` failed: %s. Stderr was:\n\n%s", err, stderr)
	}
	return vcs.CommitID(bytes.TrimSpace(stdout)), nil
}

// branchFilter is a filter for branch names.
// If not empty, only contained branch names are allowed. If empty, all names are allowed.
// The map should be made so it's not nil.
type branchFilter map[string]struct{}

// allows will return true if the current filter set-up validates against
// the passed string. If there are no filters, all strings pass.
func (f branchFilter) allows(name string) bool {
	if len(f) == 0 {
		return true
	}
	_, ok := f[name]
	return ok
}

// add adds a slice of strings to the filter.
func (f branchFilter) add(list []string) {
	for _, l := range list {
		f[l] = struct{}{}
	}
}

func (r *Repository) Branches(opt vcs.BranchesOptions) ([]*vcs.Branch, error) {
	defer r.trace(time.Now(), "Branches", opt)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	f := make(branchFilter)
	if opt.MergedInto != "" {
		b, err := r.branches("--merged", opt.MergedInto)
		if err != nil {
			return nil, err
		}
		f.add(b)
	}
	if opt.ContainsCommit != "" {
		b, err := r.branches("--contains=" + opt.ContainsCommit)
		if err != nil {
			return nil, err
		}
		f.add(b)
	}

	refs, err := r.showRef("--heads")
	if err != nil {
		return nil, err
	}

	var branches []*vcs.Branch
	for _, ref := range refs {
		name := strings.TrimPrefix(ref[1], "refs/heads/")
		id := vcs.CommitID(ref[0])
		if !f.allows(name) {
			continue
		}

		branch := &vcs.Branch{Name: name, Head: id}
		if opt.IncludeCommit {
			branch.Commit, err = r.getCommit(id)
			if err != nil {
				return nil, err
			}
		}
		if opt.BehindAheadBranch != "" {
			branch.Counts, err = r.branchesBehindAhead(name, opt.BehindAheadBranch)
			if err != nil {
				return nil, err
			}
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

// branches runs the `git branch` command followed by the given arguments and
// returns the list of branches if successful.
func (r *Repository) branches(args ...string) ([]string, error) {
	cmd := gitserver.Command("git", append([]string{"branch"}, args...)...)
	cmd.Dir = r.Dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec %v in %s failed: %v (output follows)\n\n%s", cmd.Args, cmd.Dir, err, out)
	}
	lines := strings.Split(string(out), "\n")
	lines = lines[:len(lines)-1]
	branches := make([]string, len(lines))
	for i, line := range lines {
		branches[i] = line[2:]
	}
	return branches, nil
}

// branchesBehindAhead returns the behind/ahead commit counts information for branch, against base branch.
func (r *Repository) branchesBehindAhead(branch, base string) (*vcs.BehindAhead, error) {
	if err := checkSpecArgSafety(branch); err != nil {
		return nil, err
	}
	if err := checkSpecArgSafety(base); err != nil {
		return nil, err
	}

	cmd := gitserver.Command("git", "rev-list", "--count", "--left-right", fmt.Sprintf("refs/heads/%s...refs/heads/%s", base, branch))
	cmd.Dir = r.Dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	behindAhead := strings.Split(strings.TrimSuffix(string(out), "\n"), "\t")
	b, err := strconv.ParseUint(behindAhead[0], 10, 0)
	if err != nil {
		return nil, err
	}
	a, err := strconv.ParseUint(behindAhead[1], 10, 0)
	if err != nil {
		return nil, err
	}
	return &vcs.BehindAhead{Behind: uint32(b), Ahead: uint32(a)}, nil
}

func (r *Repository) Tags() ([]*vcs.Tag, error) {
	defer r.trace(time.Now(), "Tags")

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	refs, err := r.showRef("--tags")
	if err != nil {
		return nil, err
	}

	tags := make([]*vcs.Tag, len(refs))
	for i, ref := range refs {
		tags[i] = &vcs.Tag{
			Name:     strings.TrimPrefix(ref[1], "refs/tags/"),
			CommitID: vcs.CommitID(ref[0]),
		}
	}
	return tags, nil
}

type byteSlices [][]byte

func (p byteSlices) Len() int           { return len(p) }
func (p byteSlices) Less(i, j int) bool { return bytes.Compare(p[i], p[j]) < 0 }
func (p byteSlices) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (r *Repository) showRef(arg string) ([][2]string, error) {
	cmd := gitserver.Command("git", "show-ref", arg)
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if err == vcs.ErrRepoNotExist {
			return nil, err
		}
		// Exit status of 1 and no output means there were no
		// results. This is not a fatal error.
		if cmd.ExitStatus == 1 && len(out) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("exec `git show-ref %s` in %s failed: %s. Output was:\n\n%s", arg, r.Dir, err, out)
	}

	out = bytes.TrimSuffix(out, []byte("\n")) // remove trailing newline
	lines := bytes.Split(out, []byte("\n"))
	sort.Sort(byteSlices(lines)) // sort for consistency
	refs := make([][2]string, len(lines))
	for i, line := range lines {
		if len(line) <= 41 {
			return nil, errors.New("unexpectedly short (<=41 bytes) line in `git show-ref ...` output")
		}
		id := line[:40]
		name := line[41:]
		refs[i] = [2]string{string(id), string(name)}
	}
	return refs, nil
}

// getCommit returns the commit with the given id. The caller must be holding r.editLock.
func (r *Repository) getCommit(id vcs.CommitID) (*vcs.Commit, error) {
	if err := checkSpecArgSafety(string(id)); err != nil {
		return nil, err
	}

	commits, _, err := r.commitLog(vcs.CommitsOptions{Head: id, N: 1, NoTotal: true})
	if err != nil {
		return nil, err
	}

	if len(commits) != 1 {
		return nil, fmt.Errorf("git log: expected 1 commit, got %d", len(commits))
	}

	return commits[0], nil
}

func (r *Repository) GetCommit(id vcs.CommitID) (*vcs.Commit, error) {
	defer r.trace(time.Now(), "GetCommit", id)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	return r.getCommit(id)
}

func (r *Repository) Commits(opt vcs.CommitsOptions) ([]*vcs.Commit, uint, error) {
	defer r.trace(time.Now(), "Commits", opt)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	if err := checkSpecArgSafety(string(opt.Head)); err != nil {
		return nil, 0, err
	}
	if err := checkSpecArgSafety(string(opt.Base)); err != nil {
		return nil, 0, err
	}

	return r.commitLog(opt)
}

func isBadObjectErr(output, obj string) bool {
	return string(output) == "fatal: bad object "+obj
}

func isInvalidRevisionRangeError(output, obj string) bool {
	return strings.HasPrefix(output, "fatal: Invalid revision range "+obj)
}

// commitLog returns a list of commits, and total number of commits
// starting from Head until Base or beginning of branch (unless NoTotal is true).
//
// The caller is responsible for doing checkSpecArgSafety on opt.Head and opt.Base.
func (r *Repository) commitLog(opt vcs.CommitsOptions) ([]*vcs.Commit, uint, error) {
	args := []string{"log", `--format=format:%H%x00%aN%x00%aE%x00%at%x00%cN%x00%cE%x00%ct%x00%B%x00%P%x00`}
	if opt.N != 0 {
		args = append(args, "-n", strconv.FormatUint(uint64(opt.N), 10))
	}
	if opt.Skip != 0 {
		args = append(args, "--skip="+strconv.FormatUint(uint64(opt.Skip), 10))
	}

	if opt.Path != "" {
		args = append(args, "--follow")
	}

	// Range
	rng := string(opt.Head)
	if opt.Base != "" {
		rng += "..." + string(opt.Base)
	}
	args = append(args, rng)

	if opt.Path != "" {
		args = append(args, "--", opt.Path)
	}

	cmd := gitserver.Command("git", args...)
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		out = bytes.TrimSpace(out)
		if isBadObjectErr(string(out), string(opt.Head)) {
			return nil, 0, vcs.ErrRevisionNotFound
		}
		return nil, 0, fmt.Errorf("exec `git log` failed: %s. Output was:\n\n%s", err, out)
	}

	const partsPerCommit = 9 // number of \x00-separated fields per commit
	allParts := bytes.Split(out, []byte{'\x00'})
	numCommits := len(allParts) / partsPerCommit
	commits := make([]*vcs.Commit, numCommits)
	for i := 0; i < numCommits; i++ {
		parts := allParts[partsPerCommit*i : partsPerCommit*(i+1)]

		// log outputs are newline separated, so all but the 1st commit ID part
		// has an erroneous leading newline.
		parts[0] = bytes.TrimPrefix(parts[0], []byte{'\n'})

		authorTime, err := strconv.ParseInt(string(parts[3]), 10, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("parsing git commit author time: %s", err)
		}
		committerTime, err := strconv.ParseInt(string(parts[6]), 10, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("parsing git commit committer time: %s", err)
		}

		var parents []vcs.CommitID
		if parentPart := parts[8]; len(parentPart) > 0 {
			parentIDs := bytes.Split(parentPart, []byte{' '})
			parents = make([]vcs.CommitID, len(parentIDs))
			for i, id := range parentIDs {
				parents[i] = vcs.CommitID(id)
			}
		}

		commits[i] = &vcs.Commit{
			ID:        vcs.CommitID(parts[0]),
			Author:    vcs.Signature{string(parts[1]), string(parts[2]), pbtypes.NewTimestamp(time.Unix(authorTime, 0))},
			Committer: &vcs.Signature{string(parts[4]), string(parts[5]), pbtypes.NewTimestamp(time.Unix(committerTime, 0))},
			Message:   string(bytes.TrimSuffix(parts[7], []byte{'\n'})),
			Parents:   parents,
		}
	}

	// Count commits.
	var total uint
	if !opt.NoTotal {
		cmd = gitserver.Command("git", "rev-list", "--count", rng)
		if opt.Path != "" {
			// This doesn't include --follow flag because rev-list doesn't support it, so the number may be slightly off.
			cmd.Args = append(cmd.Args, "--", opt.Path)
		}
		cmd.Dir = r.Dir
		out, err = cmd.CombinedOutput()
		if err != nil {
			return nil, 0, fmt.Errorf("exec `git rev-list --count` failed: %s. Output was:\n\n%s", err, out)
		}
		out = bytes.TrimSpace(out)
		total, err = parseUint(string(out))
		if err != nil {
			return nil, 0, err
		}
	}

	return commits, total, nil
}

func parseUint(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	return uint(n), err
}

var diffCache = synclru.New(lru.New(100))

func (r *Repository) Diff(base, head vcs.CommitID, opt *vcs.DiffOptions) (*vcs.Diff, error) {
	ensureAbsCommit(base)
	ensureAbsCommit(head)
	if opt == nil {
		opt = &vcs.DiffOptions{}
	}
	optData, err := binary.Marshal(opt)
	if err != nil {
		return nil, err
	}
	cacheKey := r.GitRootDir() + "|" + string(base) + "|" + string(head) + "|" + string(optData)
	if diff, found := diffCache.Get(cacheKey); found {
		return diff.(*vcs.Diff), nil
	}

	defer r.trace(time.Now(), "Diff", base, head, opt)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	if strings.HasPrefix(string(base), "-") || strings.HasPrefix(string(head), "-") {
		// Protect against base or head that is interpreted as command-line option.
		return nil, errors.New("diff revspecs must not start with '-'")
	}

	if opt == nil {
		opt = &vcs.DiffOptions{}
	}
	args := []string{"diff", "--full-index"}
	if opt.DetectRenames {
		args = append(args, "-M")
	}
	args = append(args, "--src-prefix="+opt.OrigPrefix)
	args = append(args, "--dst-prefix="+opt.NewPrefix)

	rng := string(base)
	if opt.ExcludeReachableFromBoth {
		rng += "..." + string(head)
	} else {
		rng += ".." + string(head)
	}

	args = append(args, rng, "--")
	cmd := gitserver.Command("git", args...)
	if opt != nil {
		cmd.Args = append(cmd.Args, opt.Paths...)
	}
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		out = bytes.TrimSpace(out)
		if isBadObjectErr(string(out), string(base)) || isBadObjectErr(string(out), string(head)) || isInvalidRevisionRangeError(string(out), string(base)) || isInvalidRevisionRangeError(string(out), string(head)) {
			return nil, vcs.ErrRevisionNotFound
		}
		return nil, fmt.Errorf("exec `git diff` failed: %s. Output was:\n\n%s", err, out)
	}
	diff := &vcs.Diff{Raw: string(out)}
	diffCache.Add(cacheKey, diff)
	return diff, nil
}

func (r *Repository) GitRootDir() string { return r.Dir }

func (r *Repository) CrossRepoDiff(base vcs.CommitID, headRepo vcs.Repository, head vcs.CommitID, opt *vcs.DiffOptions) (*vcs.Diff, error) {
	defer r.trace(time.Now(), "CrossRepoDiff", base, headRepo.GitRootDir(), head, opt)

	headDir := headRepo.GitRootDir()

	if headDir == r.Dir {
		return r.Diff(base, head, opt)
	}

	if err := r.fetchRemote(headDir); err != nil {
		return nil, err
	}

	return r.Diff(base, head, opt)
}

func (r *Repository) fetchRemote(repoDir string) error {
	r.editLock.Lock()
	defer r.editLock.Unlock()

	name := base64.URLEncoding.EncodeToString([]byte(repoDir))

	// Fetch remote commit data.
	cmd := gitserver.Command("git", "fetch", "-v", filepath.ToSlash(repoDir), "+refs/heads/*:refs/remotes/"+name+"/*")
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec %v in %s failed: %s. Output was:\n\n%s", cmd.Args, cmd.Dir, err, out)
	}
	return nil
}

// UpdateEverything updates all branches, tags, etc., to match the
// default remote repository.
func (r *Repository) UpdateEverything(opt vcs.RemoteOpts) (*vcs.UpdateResult, error) {
	defer r.trace(time.Now(), "UpdateEverything", opt)

	r.editLock.Lock()
	defer r.editLock.Unlock()

	cmd := exec.Command("git", "remote", "update", "--prune")
	cmd.Dir = r.Dir

	if opt.SSH != nil {
		gitSSHWrapper, gitSSHWrapperDir, keyFile, err := makeGitSSHWrapper(opt.SSH.PrivateKey)
		defer func() {
			if keyFile != "" {
				if err := os.Remove(keyFile); err != nil {
					log.Fatalf("Error removing SSH key file %s: %s.", keyFile, err)
				}
			}
		}()
		if err != nil {
			return nil, err
		}
		defer os.Remove(gitSSHWrapper)
		if gitSSHWrapperDir != "" {
			defer os.RemoveAll(gitSSHWrapperDir)
		}
		cmd.Env = []string{"GIT_SSH=" + gitSSHWrapper}
	}

	if opt.HTTPS != nil {
		env := environ(os.Environ())
		env.Unset("GIT_TERMINAL_PROMPT")

		gitPassHelper, gitPassHelperDir, err := makeGitPassHelper(opt.HTTPS.Pass)
		if err != nil {
			return nil, err
		}
		defer os.Remove(gitPassHelper)
		if gitPassHelperDir != "" {
			defer os.RemoveAll(gitPassHelperDir)
		}
		env = append(env, "GIT_ASKPASS="+gitPassHelper)

		cmd.Env = env
	}

	_, stderr, err := dividedOutput(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec `git remote update` failed: %v. Stderr was:\n\n%s", err, string(stderr))
	}
	result, err := parseRemoteUpdate(stderr)
	if err != nil {
		return nil, fmt.Errorf("parsing output of `git remote update` failed: %v", err)
	}
	return &result, nil
}

func (r *Repository) BlameFile(path string, opt *vcs.BlameOptions) ([]*vcs.Hunk, error) {
	defer r.trace(time.Now(), "BlameFile", path, opt)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	if opt == nil {
		opt = &vcs.BlameOptions{}
	}
	if opt.OldestCommit != "" {
		return nil, fmt.Errorf("OldestCommit not implemented")
	}
	if err := checkSpecArgSafety(string(opt.NewestCommit)); err != nil {
		return nil, err
	}
	if err := checkSpecArgSafety(string(opt.OldestCommit)); err != nil {
		return nil, err
	}

	args := []string{"blame", "-w", "--porcelain"}
	if opt.StartLine != 0 || opt.EndLine != 0 {
		args = append(args, fmt.Sprintf("-L%d,%d", opt.StartLine, opt.EndLine))
	}
	args = append(args, string(opt.NewestCommit), "--", filepath.ToSlash(path))
	cmd := gitserver.Command("git", args...)
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("exec `git blame` failed: %s. Output was:\n\n%s", err, out)
	}
	if len(out) < 1 {
		// go 1.8.5 changed the behavior of `git blame` on empty files.
		// previously, it returned a boundary commit. now, it returns nothing.
		// TODO(sqs) TODO(beyang): make `git blame` return the boundary commit
		// on an empty file somehow, or come up with some other workaround.
		st, err := os.Stat(filepath.Join(r.Dir, path))
		if err == nil && st.Size() == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("Expected git output of length at least 1")
	}

	commits := make(map[string]vcs.Commit)
	hunks := make([]*vcs.Hunk, 0)
	remainingLines := strings.Split(string(out[:len(out)-1]), "\n")
	byteOffset := 0
	for len(remainingLines) > 0 {
		// Consume hunk
		hunkHeader := strings.Split(remainingLines[0], " ")
		if len(hunkHeader) != 4 {
			fmt.Printf("Remaining lines: %+v, %d, '%s'\n", remainingLines, len(remainingLines), remainingLines[0])
			return nil, fmt.Errorf("Expected at least 4 parts to hunkHeader, but got: '%s'", hunkHeader)
		}
		commitID := hunkHeader[0]
		lineNoCur, _ := strconv.Atoi(hunkHeader[2])
		nLines, _ := strconv.Atoi(hunkHeader[3])
		hunk := &vcs.Hunk{
			CommitID:  vcs.CommitID(commitID),
			StartLine: int(lineNoCur),
			EndLine:   int(lineNoCur + nLines),
			StartByte: byteOffset,
		}

		if _, in := commits[commitID]; in {
			// Already seen commit
			byteOffset += len(remainingLines[1])
			remainingLines = remainingLines[2:]
		} else {
			// New commit
			author := strings.Join(strings.Split(remainingLines[1], " ")[1:], " ")
			email := strings.Join(strings.Split(remainingLines[2], " ")[1:], " ")
			if len(email) >= 2 && email[0] == '<' && email[len(email)-1] == '>' {
				email = email[1 : len(email)-1]
			}
			authorTime, err := strconv.ParseInt(strings.Join(strings.Split(remainingLines[3], " ")[1:], " "), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse author-time %q", remainingLines[3])
			}
			summary := strings.Join(strings.Split(remainingLines[9], " ")[1:], " ")
			commit := vcs.Commit{
				ID:      vcs.CommitID(commitID),
				Message: summary,
				Author: vcs.Signature{
					Name:  author,
					Email: email,
					Date:  pbtypes.NewTimestamp(time.Unix(authorTime, 0).In(time.UTC)),
				},
			}

			if len(remainingLines) >= 13 && strings.HasPrefix(remainingLines[10], "previous ") {
				byteOffset += len(remainingLines[12])
				remainingLines = remainingLines[13:]
			} else if len(remainingLines) >= 13 && remainingLines[10] == "boundary" {
				byteOffset += len(remainingLines[12])
				remainingLines = remainingLines[13:]
			} else if len(remainingLines) >= 12 {
				byteOffset += len(remainingLines[11])
				remainingLines = remainingLines[12:]
			} else if len(remainingLines) == 11 {
				// Empty file
				remainingLines = remainingLines[11:]
			} else {
				return nil, fmt.Errorf("Unexpected number of remaining lines (%d):\n%s", len(remainingLines), "  "+strings.Join(remainingLines, "\n  "))
			}

			commits[commitID] = commit
		}

		if commit, present := commits[commitID]; present {
			// Should always be present, but check just to avoid
			// panicking in case of a (somewhat likely) bug in our
			// git-blame parser above.
			hunk.CommitID = commit.ID
			hunk.Author = commit.Author
		}

		// Consume remaining lines in hunk
		for i := 1; i < nLines; i++ {
			byteOffset += len(remainingLines[1])
			remainingLines = remainingLines[2:]
		}

		hunk.EndByte = byteOffset
		hunks = append(hunks, hunk)
	}

	return hunks, nil
}

func (r *Repository) MergeBase(a, b vcs.CommitID) (vcs.CommitID, error) {
	defer r.trace(time.Now(), "MergeBase", a, b)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	cmd := gitserver.Command("git", "merge-base", "--", string(a), string(b))
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec %v failed: %s. Output was:\n\n%s", cmd.Args, err, out)
	}
	return vcs.CommitID(bytes.TrimSpace(out)), nil
}

func (r *Repository) CrossRepoMergeBase(a vcs.CommitID, repoB vcs.Repository, b vcs.CommitID) (vcs.CommitID, error) {
	defer r.trace(time.Now(), "CrossRepoMergeBase", a, repoB.GitRootDir(), b)

	// git.Repository inherits GitRootDir and CrossRepo from its
	// embedded gitcmd.Repository.

	repoBDir := repoB.GitRootDir()

	if repoBDir != r.Dir {
		if err := r.fetchRemote(repoBDir); err != nil {
			return "", err
		}
	}

	return r.MergeBase(a, b)
}

func (r *Repository) Search(at vcs.CommitID, opt vcs.SearchOptions) ([]*vcs.SearchResult, error) {
	defer r.trace(time.Now(), "Search", at, opt)

	return gitserver.Search(r.Dir, at, opt)
}

func (r *Repository) Committers(opt vcs.CommittersOptions) ([]*vcs.Committer, error) {
	defer r.trace(time.Now(), "Committers", opt)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	if opt.Rev == "" {
		opt.Rev = "HEAD"
	}

	cmd := gitserver.Command("git", "shortlog", "-sne", opt.Rev)
	cmd.Dir = r.Dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec `git shortlog -sne` failed: %v", err)
	}
	out = bytes.TrimSpace(out)

	allEntries := bytes.Split(out, []byte{'\n'})
	numEntries := len(allEntries)
	if opt.N > 0 && numEntries > opt.N {
		numEntries = opt.N
	}
	var committers []*vcs.Committer
	for i := 0; i < numEntries; i++ {
		line := string(allEntries[i])
		if match := logEntryPattern.FindStringSubmatch(line); match != nil {
			commits, err2 := strconv.Atoi(match[1])
			if err2 != nil {
				continue
			}
			committers = append(committers, &vcs.Committer{
				Commits: int32(commits),
				Name:    match[2],
				Email:   match[3],
			})
		}
	}
	return committers, nil
}

func (r *Repository) ReadFile(commit vcs.CommitID, name string) ([]byte, error) {
	defer r.trace(time.Now(), "ReadFile", name)

	if err := checkSpecArgSafety(string(commit)); err != nil {
		return nil, err
	}

	name = internal.Rel(name)
	r.editLock.RLock()
	defer r.editLock.RUnlock()
	b, err := r.readFileBytes(commit, name)
	if err != nil {
		return nil, err
	}
	return b, nil
}

var readFileBytesCache = synclru.New(lru.New(1000))

func (r *Repository) readFileBytes(commit vcs.CommitID, name string) ([]byte, error) {
	ensureAbsCommit(commit)
	cacheKey := r.GitRootDir() + "|" + string(commit) + "|" + name
	if data, found := readFileBytesCache.Get(cacheKey); found {
		return data.([]byte), nil
	}

	cmd := gitserver.Command("git", "show", string(commit)+":"+name)
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if bytes.Contains(out, []byte("exists on disk, but not in")) || bytes.Contains(out, []byte("does not exist")) {
			return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
		}
		if bytes.HasPrefix(out, []byte("fatal: bad object ")) {
			// Could be a git submodule.
			fi, err := r.Stat(commit, name)
			if err != nil {
				return nil, err
			}
			// Return empty for a submodule for now.
			if fi.Mode()&vcs.ModeSubmodule != 0 {
				return nil, nil
			}

		}
		return nil, fmt.Errorf("exec %v failed: %s. Output was:\n\n%s", cmd.Args, err, out)
	}
	readFileBytesCache.Add(cacheKey, out)
	return out, nil
}

func (r *Repository) Lstat(commit vcs.CommitID, path string) (os.FileInfo, error) {
	defer r.trace(time.Now(), "Lstat", path)

	if err := checkSpecArgSafety(string(commit)); err != nil {
		return nil, err
	}

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	path = filepath.Clean(internal.Rel(path))

	if path == "." {
		// Special case root, which is not returned by `git ls-tree`.
		return &util.FileInfo{Mode_: os.ModeDir}, nil
	}

	fis, err := r.lsTree(commit, path, false)
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return nil, &os.PathError{Op: "ls-tree", Path: path, Err: os.ErrNotExist}
	}

	return fis[0], nil
}

func (r *Repository) Stat(commit vcs.CommitID, path string) (os.FileInfo, error) {
	defer r.trace(time.Now(), "Stat", path)

	if err := checkSpecArgSafety(string(commit)); err != nil {
		return nil, err
	}

	path = internal.Rel(path)

	r.editLock.RLock()
	defer r.editLock.RUnlock()

	fi, err := r.Lstat(commit, path)
	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		// Deref symlink.
		b, err := r.readFileBytes(commit, path)
		if err != nil {
			return nil, err
		}
		fi2, err := r.Lstat(commit, string(b))
		if err != nil {
			return nil, err
		}
		fi2.(*util.FileInfo).Name_ = fi.Name()
		return fi2, nil
	}

	return fi, nil
}

func (r *Repository) ReadDir(commit vcs.CommitID, path string, recurse bool) ([]os.FileInfo, error) {
	defer r.trace(time.Now(), "ReadDir", path)

	if err := checkSpecArgSafety(string(commit)); err != nil {
		return nil, err
	}

	r.editLock.RLock()
	defer r.editLock.RUnlock()
	// Trailing slash is necessary to ls-tree under the dir (not just
	// to list the dir's tree entry in its parent dir).
	return r.lsTree(commit, filepath.Clean(internal.Rel(path))+"/", recurse)
}

var lsTreeCache = synclru.New(lru.New(10000))

// lsTree returns ls of tree at path. The caller must be holding r.editLock.RLock().
func (r *Repository) lsTree(commit vcs.CommitID, path string, recurse bool) ([]os.FileInfo, error) {
	ensureAbsCommit(commit)
	cacheKey := r.GitRootDir() + "|" + string(commit) + "|" + path + "|" + strconv.FormatBool(recurse)
	if fis, found := lsTreeCache.Get(cacheKey); found {
		return fis.([]os.FileInfo), nil
	}

	// Don't call filepath.Clean(path) because ReadDir needs to pass
	// path with a trailing slash.

	if err := checkSpecArgSafety(path); err != nil {
		return nil, err
	}

	args := []string{"ls-tree", "--full-name", string(commit)}
	if recurse {
		args = append(args, "-r", "-t")
	}
	args = append(args, "--", filepath.ToSlash(path))
	cmd := gitserver.Command("git", args...)
	cmd.Dir = r.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if bytes.Contains(out, []byte("exists on disk, but not in")) {
			return nil, &os.PathError{Op: "ls-tree", Path: filepath.ToSlash(path), Err: os.ErrNotExist}
		}
		return nil, fmt.Errorf("exec %v failed: %s. Output was:\n\n%s", cmd.Args, err, out)
	}

	if len(out) == 0 {
		return nil, os.ErrNotExist
	}

	prefixLen := strings.LastIndexByte(strings.TrimPrefix(path, "./"), '/') + 1
	lines := strings.Split(string(out), "\n")
	fis := make([]os.FileInfo, len(lines)-1)
	for i, line := range lines {
		if i == len(lines)-1 {
			// last entry is empty
			continue
		}

		tabPos := strings.IndexByte(line, '\t')
		if tabPos == -1 {
			return nil, fmt.Errorf("invalid `git ls-tree` output: %q", out)
		}
		info := strings.SplitN(line[:tabPos], " ", 3)
		name := line[tabPos+1:]

		if len(info) != 3 {
			return nil, fmt.Errorf("invalid `git ls-tree` output: %q", out)
		}
		typ := info[1]
		oid := info[2]
		if len(oid) != 40 {
			return nil, fmt.Errorf("invalid `git ls-tree` oid output: %q", oid)
		}

		var sys interface{}
		mode, err := strconv.ParseInt(info[0], 8, 32)
		if err != nil {
			return nil, err
		}
		switch typ {
		case "blob":
			const gitModeSymlink = 020000
			if mode&gitModeSymlink != 0 {
				mode = int64(os.ModeSymlink)
			} else {
				// Regular file.
				mode = mode | 0644
			}
		case "commit":
			mode = mode | vcs.ModeSubmodule
			cmd := gitserver.Command("git", "config", "--get", "submodule."+name+".url")
			cmd.Dir = r.Dir
			url := "" // url is not available if submodules are not initialized
			if out, err := cmd.Output(); err == nil {
				url = string(bytes.TrimSpace(out))
			}
			sys = vcs.SubmoduleInfo{
				URL:      url,
				CommitID: vcs.CommitID(oid),
			}
		case "tree":
			mode = mode | int64(os.ModeDir)
		}

		fis[i] = &util.FileInfo{
			Name_: name[prefixLen:],
			Mode_: os.FileMode(mode),
			Sys_:  sys,
		}
	}
	util.SortFileInfosByName(fis)

	lsTreeCache.Add(cacheKey, fis)
	return fis, nil
}

// makeGitSSHWrapper writes a GIT_SSH wrapper that runs ssh with the
// private key. You should remove the sshWrapper, sshWrapperDir and
// the keyFile after using them.
func makeGitSSHWrapper(privKey []byte) (sshWrapper, sshWrapperDir, keyFile string, err error) {
	var otherOpt string
	if InsecureSkipCheckVerifySSH {
		otherOpt = "-o StrictHostKeyChecking=no"
	}

	kf, err := ioutil.TempFile("", "go-vcs-gitcmd-key")
	if err != nil {
		return "", "", "", err
	}
	keyFile = kf.Name()
	err = internal.WriteFileWithPermissions(keyFile, privKey, 0600)
	if err != nil {
		return "", "", keyFile, err
	}

	tmpFile, tmpFileDir, err := gitSSHWrapper(keyFile, otherOpt)
	return tmpFile, tmpFileDir, keyFile, err
}

// makeGitPassHelper writes a GIT_ASKPASS helper that supplies password over stdout.
// You should remove the passHelper (and tempDir if any) after using it.
func makeGitPassHelper(pass string) (passHelper string, tempDir string, err error) {

	tmpFile, dir, err := internal.ScriptFile("go-vcs-gitcmd-ask")
	if err != nil {
		return tmpFile, dir, err
	}

	passPath := filepath.Join(dir, "password")
	err = internal.WriteFileWithPermissions(passPath, []byte(pass), 0600)
	if err != nil {
		return tmpFile, dir, err
	}

	var script string

	// We assume passPath can be escaped with a simple wrapping of single
	// quotes. The path is not user controlled so this assumption should
	// not be violated.
	if runtime.GOOS == "windows" {
		script = "@echo off\ntype " + passPath + "\n"
	} else {
		script = "#!/bin/sh\ncat '" + passPath + "'\n"
	}

	err = internal.WriteFileWithPermissions(tmpFile, []byte(script), 0500)
	return tmpFile, dir, err
}

// InsecureSkipCheckVerifySSH controls whether the client verifies the
// SSH server's certificate or host key. If InsecureSkipCheckVerifySSH
// is true, the program is susceptible to a man-in-the-middle
// attack. This should only be used for testing.
var InsecureSkipCheckVerifySSH bool

// environ is a slice of strings representing the environment, in the form "key=value".
type environ []string

// Unset a single environment variable.
func (e *environ) Unset(key string) {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			(*e)[i] = (*e)[len(*e)-1]
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

// Makes system-dependent SSH wrapper
func gitSSHWrapper(keyFile string, otherOpt string) (sshWrapperFile string, tempDir string, err error) {
	// TODO(sqs): encrypt and store the key in the env so that
	// attackers can't decrypt if they have disk access after our
	// process dies

	var script string

	if runtime.GOOS == "windows" {
		script = `
	@echo off
	ssh -o ControlMaster=no -o ControlPath=none ` + otherOpt + ` -i ` + filepath.ToSlash(keyFile) + ` "%@"
`
	} else {
		script = `
	#!/bin/sh
	exec /usr/bin/ssh -o ControlMaster=no -o ControlPath=none ` + otherOpt + ` -i ` + filepath.ToSlash(keyFile) + ` "$@"
`
	}

	sshWrapperName, tempDir, err := internal.ScriptFile("go-vcs-gitcmd")
	if err != nil {
		return sshWrapperName, tempDir, err
	}

	err = internal.WriteFileWithPermissions(sshWrapperName, []byte(script), 0500)
	return sshWrapperName, tempDir, err
}

func ensureAbsCommit(commitID vcs.CommitID) {
	// We don't want to even be running commands on non-absolute
	// commit IDs if we can avoid it, because we can't cache the
	// expensive part of those computations.
	if len(commitID) != 40 {
		panic("non-absolute commit ID")
	}
}
