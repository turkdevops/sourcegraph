// generated by gen-mocks; DO NOT EDIT

package mockstore

import (
	"golang.org/x/net/context"
	"src.sourcegraph.com/sourcegraph/store"
)

type RepoPerms struct {
	ListRepoUsers_ func(ctx context.Context, repo string) ([]int32, error)
}

func (s *RepoPerms) ListRepoUsers(ctx context.Context, repo string) ([]int32, error) {
	return s.ListRepoUsers_(ctx, repo)
}

var _ store.RepoPerms = (*RepoPerms)(nil)
