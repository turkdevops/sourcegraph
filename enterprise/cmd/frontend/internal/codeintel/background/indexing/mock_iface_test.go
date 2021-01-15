// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.

package indexing

import (
	"context"
	dbstore "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/dbstore"
	"regexp"
	"sync"
	"time"
)

// MockDBStore is a mock implementation of the DBStore interface (from the
// package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/codeintel/background/indexing)
// used for unit testing.
type MockDBStore struct {
	// GetRepositoriesWithIndexConfigurationFunc is an instance of a mock
	// function object controlling the behavior of the method
	// GetRepositoriesWithIndexConfiguration.
	GetRepositoriesWithIndexConfigurationFunc *DBStoreGetRepositoriesWithIndexConfigurationFunc
	// IndexableRepositoriesFunc is an instance of a mock function object
	// controlling the behavior of the method IndexableRepositories.
	IndexableRepositoriesFunc *DBStoreIndexableRepositoriesFunc
	// RepoUsageStatisticsFunc is an instance of a mock function object
	// controlling the behavior of the method RepoUsageStatistics.
	RepoUsageStatisticsFunc *DBStoreRepoUsageStatisticsFunc
	// ResetIndexableRepositoriesFunc is an instance of a mock function
	// object controlling the behavior of the method
	// ResetIndexableRepositories.
	ResetIndexableRepositoriesFunc *DBStoreResetIndexableRepositoriesFunc
	// UpdateIndexableRepositoryFunc is an instance of a mock function
	// object controlling the behavior of the method
	// UpdateIndexableRepository.
	UpdateIndexableRepositoryFunc *DBStoreUpdateIndexableRepositoryFunc
}

// NewMockDBStore creates a new mock of the DBStore interface. All methods
// return zero values for all results, unless overwritten.
func NewMockDBStore() *MockDBStore {
	return &MockDBStore{
		GetRepositoriesWithIndexConfigurationFunc: &DBStoreGetRepositoriesWithIndexConfigurationFunc{
			defaultHook: func(context.Context) ([]int, error) {
				return nil, nil
			},
		},
		IndexableRepositoriesFunc: &DBStoreIndexableRepositoriesFunc{
			defaultHook: func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error) {
				return nil, nil
			},
		},
		RepoUsageStatisticsFunc: &DBStoreRepoUsageStatisticsFunc{
			defaultHook: func(context.Context) ([]dbstore.RepoUsageStatistics, error) {
				return nil, nil
			},
		},
		ResetIndexableRepositoriesFunc: &DBStoreResetIndexableRepositoriesFunc{
			defaultHook: func(context.Context, time.Time) error {
				return nil
			},
		},
		UpdateIndexableRepositoryFunc: &DBStoreUpdateIndexableRepositoryFunc{
			defaultHook: func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error {
				return nil
			},
		},
	}
}

// NewMockDBStoreFrom creates a new mock of the MockDBStore interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockDBStoreFrom(i DBStore) *MockDBStore {
	return &MockDBStore{
		GetRepositoriesWithIndexConfigurationFunc: &DBStoreGetRepositoriesWithIndexConfigurationFunc{
			defaultHook: i.GetRepositoriesWithIndexConfiguration,
		},
		IndexableRepositoriesFunc: &DBStoreIndexableRepositoriesFunc{
			defaultHook: i.IndexableRepositories,
		},
		RepoUsageStatisticsFunc: &DBStoreRepoUsageStatisticsFunc{
			defaultHook: i.RepoUsageStatistics,
		},
		ResetIndexableRepositoriesFunc: &DBStoreResetIndexableRepositoriesFunc{
			defaultHook: i.ResetIndexableRepositories,
		},
		UpdateIndexableRepositoryFunc: &DBStoreUpdateIndexableRepositoryFunc{
			defaultHook: i.UpdateIndexableRepository,
		},
	}
}

// DBStoreGetRepositoriesWithIndexConfigurationFunc describes the behavior
// when the GetRepositoriesWithIndexConfiguration method of the parent
// MockDBStore instance is invoked.
type DBStoreGetRepositoriesWithIndexConfigurationFunc struct {
	defaultHook func(context.Context) ([]int, error)
	hooks       []func(context.Context) ([]int, error)
	history     []DBStoreGetRepositoriesWithIndexConfigurationFuncCall
	mutex       sync.Mutex
}

// GetRepositoriesWithIndexConfiguration delegates to the next hook function
// in the queue and stores the parameter and result values of this
// invocation.
func (m *MockDBStore) GetRepositoriesWithIndexConfiguration(v0 context.Context) ([]int, error) {
	r0, r1 := m.GetRepositoriesWithIndexConfigurationFunc.nextHook()(v0)
	m.GetRepositoriesWithIndexConfigurationFunc.appendCall(DBStoreGetRepositoriesWithIndexConfigurationFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// GetRepositoriesWithIndexConfiguration method of the parent MockDBStore
// instance is invoked and the hook queue is empty.
func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) SetDefaultHook(hook func(context.Context) ([]int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetRepositoriesWithIndexConfiguration method of the parent MockDBStore
// instance inovkes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) PushHook(hook func(context.Context) ([]int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) SetDefaultReturn(r0 []int, r1 error) {
	f.SetDefaultHook(func(context.Context) ([]int, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) PushReturn(r0 []int, r1 error) {
	f.PushHook(func(context.Context) ([]int, error) {
		return r0, r1
	})
}

func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) nextHook() func(context.Context) ([]int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) appendCall(r0 DBStoreGetRepositoriesWithIndexConfigurationFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// DBStoreGetRepositoriesWithIndexConfigurationFuncCall objects describing
// the invocations of this function.
func (f *DBStoreGetRepositoriesWithIndexConfigurationFunc) History() []DBStoreGetRepositoriesWithIndexConfigurationFuncCall {
	f.mutex.Lock()
	history := make([]DBStoreGetRepositoriesWithIndexConfigurationFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreGetRepositoriesWithIndexConfigurationFuncCall is an object that
// describes an invocation of method GetRepositoriesWithIndexConfiguration
// on an instance of MockDBStore.
type DBStoreGetRepositoriesWithIndexConfigurationFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []int
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreGetRepositoriesWithIndexConfigurationFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreGetRepositoriesWithIndexConfigurationFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// DBStoreIndexableRepositoriesFunc describes the behavior when the
// IndexableRepositories method of the parent MockDBStore instance is
// invoked.
type DBStoreIndexableRepositoriesFunc struct {
	defaultHook func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error)
	hooks       []func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error)
	history     []DBStoreIndexableRepositoriesFuncCall
	mutex       sync.Mutex
}

// IndexableRepositories delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockDBStore) IndexableRepositories(v0 context.Context, v1 dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error) {
	r0, r1 := m.IndexableRepositoriesFunc.nextHook()(v0, v1)
	m.IndexableRepositoriesFunc.appendCall(DBStoreIndexableRepositoriesFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// IndexableRepositories method of the parent MockDBStore instance is
// invoked and the hook queue is empty.
func (f *DBStoreIndexableRepositoriesFunc) SetDefaultHook(hook func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// IndexableRepositories method of the parent MockDBStore instance inovkes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *DBStoreIndexableRepositoriesFunc) PushHook(hook func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *DBStoreIndexableRepositoriesFunc) SetDefaultReturn(r0 []dbstore.IndexableRepository, r1 error) {
	f.SetDefaultHook(func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *DBStoreIndexableRepositoriesFunc) PushReturn(r0 []dbstore.IndexableRepository, r1 error) {
	f.PushHook(func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error) {
		return r0, r1
	})
}

func (f *DBStoreIndexableRepositoriesFunc) nextHook() func(context.Context, dbstore.IndexableRepositoryQueryOptions) ([]dbstore.IndexableRepository, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreIndexableRepositoriesFunc) appendCall(r0 DBStoreIndexableRepositoriesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreIndexableRepositoriesFuncCall
// objects describing the invocations of this function.
func (f *DBStoreIndexableRepositoriesFunc) History() []DBStoreIndexableRepositoriesFuncCall {
	f.mutex.Lock()
	history := make([]DBStoreIndexableRepositoriesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreIndexableRepositoriesFuncCall is an object that describes an
// invocation of method IndexableRepositories on an instance of MockDBStore.
type DBStoreIndexableRepositoriesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 dbstore.IndexableRepositoryQueryOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []dbstore.IndexableRepository
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreIndexableRepositoriesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreIndexableRepositoriesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// DBStoreRepoUsageStatisticsFunc describes the behavior when the
// RepoUsageStatistics method of the parent MockDBStore instance is invoked.
type DBStoreRepoUsageStatisticsFunc struct {
	defaultHook func(context.Context) ([]dbstore.RepoUsageStatistics, error)
	hooks       []func(context.Context) ([]dbstore.RepoUsageStatistics, error)
	history     []DBStoreRepoUsageStatisticsFuncCall
	mutex       sync.Mutex
}

// RepoUsageStatistics delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockDBStore) RepoUsageStatistics(v0 context.Context) ([]dbstore.RepoUsageStatistics, error) {
	r0, r1 := m.RepoUsageStatisticsFunc.nextHook()(v0)
	m.RepoUsageStatisticsFunc.appendCall(DBStoreRepoUsageStatisticsFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the RepoUsageStatistics
// method of the parent MockDBStore instance is invoked and the hook queue
// is empty.
func (f *DBStoreRepoUsageStatisticsFunc) SetDefaultHook(hook func(context.Context) ([]dbstore.RepoUsageStatistics, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// RepoUsageStatistics method of the parent MockDBStore instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *DBStoreRepoUsageStatisticsFunc) PushHook(hook func(context.Context) ([]dbstore.RepoUsageStatistics, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *DBStoreRepoUsageStatisticsFunc) SetDefaultReturn(r0 []dbstore.RepoUsageStatistics, r1 error) {
	f.SetDefaultHook(func(context.Context) ([]dbstore.RepoUsageStatistics, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *DBStoreRepoUsageStatisticsFunc) PushReturn(r0 []dbstore.RepoUsageStatistics, r1 error) {
	f.PushHook(func(context.Context) ([]dbstore.RepoUsageStatistics, error) {
		return r0, r1
	})
}

func (f *DBStoreRepoUsageStatisticsFunc) nextHook() func(context.Context) ([]dbstore.RepoUsageStatistics, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreRepoUsageStatisticsFunc) appendCall(r0 DBStoreRepoUsageStatisticsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreRepoUsageStatisticsFuncCall objects
// describing the invocations of this function.
func (f *DBStoreRepoUsageStatisticsFunc) History() []DBStoreRepoUsageStatisticsFuncCall {
	f.mutex.Lock()
	history := make([]DBStoreRepoUsageStatisticsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreRepoUsageStatisticsFuncCall is an object that describes an
// invocation of method RepoUsageStatistics on an instance of MockDBStore.
type DBStoreRepoUsageStatisticsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []dbstore.RepoUsageStatistics
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreRepoUsageStatisticsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreRepoUsageStatisticsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// DBStoreResetIndexableRepositoriesFunc describes the behavior when the
// ResetIndexableRepositories method of the parent MockDBStore instance is
// invoked.
type DBStoreResetIndexableRepositoriesFunc struct {
	defaultHook func(context.Context, time.Time) error
	hooks       []func(context.Context, time.Time) error
	history     []DBStoreResetIndexableRepositoriesFuncCall
	mutex       sync.Mutex
}

// ResetIndexableRepositories delegates to the next hook function in the
// queue and stores the parameter and result values of this invocation.
func (m *MockDBStore) ResetIndexableRepositories(v0 context.Context, v1 time.Time) error {
	r0 := m.ResetIndexableRepositoriesFunc.nextHook()(v0, v1)
	m.ResetIndexableRepositoriesFunc.appendCall(DBStoreResetIndexableRepositoriesFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the
// ResetIndexableRepositories method of the parent MockDBStore instance is
// invoked and the hook queue is empty.
func (f *DBStoreResetIndexableRepositoriesFunc) SetDefaultHook(hook func(context.Context, time.Time) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// ResetIndexableRepositories method of the parent MockDBStore instance
// inovkes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *DBStoreResetIndexableRepositoriesFunc) PushHook(hook func(context.Context, time.Time) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *DBStoreResetIndexableRepositoriesFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, time.Time) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *DBStoreResetIndexableRepositoriesFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, time.Time) error {
		return r0
	})
}

func (f *DBStoreResetIndexableRepositoriesFunc) nextHook() func(context.Context, time.Time) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreResetIndexableRepositoriesFunc) appendCall(r0 DBStoreResetIndexableRepositoriesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreResetIndexableRepositoriesFuncCall
// objects describing the invocations of this function.
func (f *DBStoreResetIndexableRepositoriesFunc) History() []DBStoreResetIndexableRepositoriesFuncCall {
	f.mutex.Lock()
	history := make([]DBStoreResetIndexableRepositoriesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreResetIndexableRepositoriesFuncCall is an object that describes an
// invocation of method ResetIndexableRepositories on an instance of
// MockDBStore.
type DBStoreResetIndexableRepositoriesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 time.Time
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreResetIndexableRepositoriesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreResetIndexableRepositoriesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// DBStoreUpdateIndexableRepositoryFunc describes the behavior when the
// UpdateIndexableRepository method of the parent MockDBStore instance is
// invoked.
type DBStoreUpdateIndexableRepositoryFunc struct {
	defaultHook func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error
	hooks       []func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error
	history     []DBStoreUpdateIndexableRepositoryFuncCall
	mutex       sync.Mutex
}

// UpdateIndexableRepository delegates to the next hook function in the
// queue and stores the parameter and result values of this invocation.
func (m *MockDBStore) UpdateIndexableRepository(v0 context.Context, v1 dbstore.UpdateableIndexableRepository, v2 time.Time) error {
	r0 := m.UpdateIndexableRepositoryFunc.nextHook()(v0, v1, v2)
	m.UpdateIndexableRepositoryFunc.appendCall(DBStoreUpdateIndexableRepositoryFuncCall{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the
// UpdateIndexableRepository method of the parent MockDBStore instance is
// invoked and the hook queue is empty.
func (f *DBStoreUpdateIndexableRepositoryFunc) SetDefaultHook(hook func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UpdateIndexableRepository method of the parent MockDBStore instance
// inovkes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *DBStoreUpdateIndexableRepositoryFunc) PushHook(hook func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *DBStoreUpdateIndexableRepositoryFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *DBStoreUpdateIndexableRepositoryFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error {
		return r0
	})
}

func (f *DBStoreUpdateIndexableRepositoryFunc) nextHook() func(context.Context, dbstore.UpdateableIndexableRepository, time.Time) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreUpdateIndexableRepositoryFunc) appendCall(r0 DBStoreUpdateIndexableRepositoryFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreUpdateIndexableRepositoryFuncCall
// objects describing the invocations of this function.
func (f *DBStoreUpdateIndexableRepositoryFunc) History() []DBStoreUpdateIndexableRepositoryFuncCall {
	f.mutex.Lock()
	history := make([]DBStoreUpdateIndexableRepositoryFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreUpdateIndexableRepositoryFuncCall is an object that describes an
// invocation of method UpdateIndexableRepository on an instance of
// MockDBStore.
type DBStoreUpdateIndexableRepositoryFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 dbstore.UpdateableIndexableRepository
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 time.Time
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreUpdateIndexableRepositoryFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreUpdateIndexableRepositoryFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// MockGitserverClient is a mock implementation of the GitserverClient
// interface (from the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/codeintel/background/indexing)
// used for unit testing.
type MockGitserverClient struct {
	// HeadFunc is an instance of a mock function object controlling the
	// behavior of the method Head.
	HeadFunc *GitserverClientHeadFunc
	// ListFilesFunc is an instance of a mock function object controlling
	// the behavior of the method ListFiles.
	ListFilesFunc *GitserverClientListFilesFunc
}

// NewMockGitserverClient creates a new mock of the GitserverClient
// interface. All methods return zero values for all results, unless
// overwritten.
func NewMockGitserverClient() *MockGitserverClient {
	return &MockGitserverClient{
		HeadFunc: &GitserverClientHeadFunc{
			defaultHook: func(context.Context, int) (string, error) {
				return "", nil
			},
		},
		ListFilesFunc: &GitserverClientListFilesFunc{
			defaultHook: func(context.Context, int, string, *regexp.Regexp) ([]string, error) {
				return nil, nil
			},
		},
	}
}

// NewMockGitserverClientFrom creates a new mock of the MockGitserverClient
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockGitserverClientFrom(i GitserverClient) *MockGitserverClient {
	return &MockGitserverClient{
		HeadFunc: &GitserverClientHeadFunc{
			defaultHook: i.Head,
		},
		ListFilesFunc: &GitserverClientListFilesFunc{
			defaultHook: i.ListFiles,
		},
	}
}

// GitserverClientHeadFunc describes the behavior when the Head method of
// the parent MockGitserverClient instance is invoked.
type GitserverClientHeadFunc struct {
	defaultHook func(context.Context, int) (string, error)
	hooks       []func(context.Context, int) (string, error)
	history     []GitserverClientHeadFuncCall
	mutex       sync.Mutex
}

// Head delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockGitserverClient) Head(v0 context.Context, v1 int) (string, error) {
	r0, r1 := m.HeadFunc.nextHook()(v0, v1)
	m.HeadFunc.appendCall(GitserverClientHeadFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Head method of the
// parent MockGitserverClient instance is invoked and the hook queue is
// empty.
func (f *GitserverClientHeadFunc) SetDefaultHook(hook func(context.Context, int) (string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Head method of the parent MockGitserverClient instance inovkes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *GitserverClientHeadFunc) PushHook(hook func(context.Context, int) (string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *GitserverClientHeadFunc) SetDefaultReturn(r0 string, r1 error) {
	f.SetDefaultHook(func(context.Context, int) (string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *GitserverClientHeadFunc) PushReturn(r0 string, r1 error) {
	f.PushHook(func(context.Context, int) (string, error) {
		return r0, r1
	})
}

func (f *GitserverClientHeadFunc) nextHook() func(context.Context, int) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GitserverClientHeadFunc) appendCall(r0 GitserverClientHeadFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of GitserverClientHeadFuncCall objects
// describing the invocations of this function.
func (f *GitserverClientHeadFunc) History() []GitserverClientHeadFuncCall {
	f.mutex.Lock()
	history := make([]GitserverClientHeadFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GitserverClientHeadFuncCall is an object that describes an invocation of
// method Head on an instance of MockGitserverClient.
type GitserverClientHeadFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c GitserverClientHeadFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GitserverClientHeadFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// GitserverClientListFilesFunc describes the behavior when the ListFiles
// method of the parent MockGitserverClient instance is invoked.
type GitserverClientListFilesFunc struct {
	defaultHook func(context.Context, int, string, *regexp.Regexp) ([]string, error)
	hooks       []func(context.Context, int, string, *regexp.Regexp) ([]string, error)
	history     []GitserverClientListFilesFuncCall
	mutex       sync.Mutex
}

// ListFiles delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockGitserverClient) ListFiles(v0 context.Context, v1 int, v2 string, v3 *regexp.Regexp) ([]string, error) {
	r0, r1 := m.ListFilesFunc.nextHook()(v0, v1, v2, v3)
	m.ListFilesFunc.appendCall(GitserverClientListFilesFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the ListFiles method of
// the parent MockGitserverClient instance is invoked and the hook queue is
// empty.
func (f *GitserverClientListFilesFunc) SetDefaultHook(hook func(context.Context, int, string, *regexp.Regexp) ([]string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// ListFiles method of the parent MockGitserverClient instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *GitserverClientListFilesFunc) PushHook(hook func(context.Context, int, string, *regexp.Regexp) ([]string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *GitserverClientListFilesFunc) SetDefaultReturn(r0 []string, r1 error) {
	f.SetDefaultHook(func(context.Context, int, string, *regexp.Regexp) ([]string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *GitserverClientListFilesFunc) PushReturn(r0 []string, r1 error) {
	f.PushHook(func(context.Context, int, string, *regexp.Regexp) ([]string, error) {
		return r0, r1
	})
}

func (f *GitserverClientListFilesFunc) nextHook() func(context.Context, int, string, *regexp.Regexp) ([]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GitserverClientListFilesFunc) appendCall(r0 GitserverClientListFilesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of GitserverClientListFilesFuncCall objects
// describing the invocations of this function.
func (f *GitserverClientListFilesFunc) History() []GitserverClientListFilesFuncCall {
	f.mutex.Lock()
	history := make([]GitserverClientListFilesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GitserverClientListFilesFuncCall is an object that describes an
// invocation of method ListFiles on an instance of MockGitserverClient.
type GitserverClientListFilesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 *regexp.Regexp
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c GitserverClientListFilesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GitserverClientListFilesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// MockIndexEnqueuer is a mock implementation of the IndexEnqueuer interface
// (from the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/codeintel/background/indexing)
// used for unit testing.
type MockIndexEnqueuer struct {
	// QueueIndexFunc is an instance of a mock function object controlling
	// the behavior of the method QueueIndex.
	QueueIndexFunc *IndexEnqueuerQueueIndexFunc
}

// NewMockIndexEnqueuer creates a new mock of the IndexEnqueuer interface.
// All methods return zero values for all results, unless overwritten.
func NewMockIndexEnqueuer() *MockIndexEnqueuer {
	return &MockIndexEnqueuer{
		QueueIndexFunc: &IndexEnqueuerQueueIndexFunc{
			defaultHook: func(context.Context, int) error {
				return nil
			},
		},
	}
}

// NewMockIndexEnqueuerFrom creates a new mock of the MockIndexEnqueuer
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockIndexEnqueuerFrom(i IndexEnqueuer) *MockIndexEnqueuer {
	return &MockIndexEnqueuer{
		QueueIndexFunc: &IndexEnqueuerQueueIndexFunc{
			defaultHook: i.QueueIndex,
		},
	}
}

// IndexEnqueuerQueueIndexFunc describes the behavior when the QueueIndex
// method of the parent MockIndexEnqueuer instance is invoked.
type IndexEnqueuerQueueIndexFunc struct {
	defaultHook func(context.Context, int) error
	hooks       []func(context.Context, int) error
	history     []IndexEnqueuerQueueIndexFuncCall
	mutex       sync.Mutex
}

// QueueIndex delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockIndexEnqueuer) QueueIndex(v0 context.Context, v1 int) error {
	r0 := m.QueueIndexFunc.nextHook()(v0, v1)
	m.QueueIndexFunc.appendCall(IndexEnqueuerQueueIndexFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the QueueIndex method of
// the parent MockIndexEnqueuer instance is invoked and the hook queue is
// empty.
func (f *IndexEnqueuerQueueIndexFunc) SetDefaultHook(hook func(context.Context, int) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// QueueIndex method of the parent MockIndexEnqueuer instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *IndexEnqueuerQueueIndexFunc) PushHook(hook func(context.Context, int) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *IndexEnqueuerQueueIndexFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, int) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *IndexEnqueuerQueueIndexFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, int) error {
		return r0
	})
}

func (f *IndexEnqueuerQueueIndexFunc) nextHook() func(context.Context, int) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *IndexEnqueuerQueueIndexFunc) appendCall(r0 IndexEnqueuerQueueIndexFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of IndexEnqueuerQueueIndexFuncCall objects
// describing the invocations of this function.
func (f *IndexEnqueuerQueueIndexFunc) History() []IndexEnqueuerQueueIndexFuncCall {
	f.mutex.Lock()
	history := make([]IndexEnqueuerQueueIndexFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// IndexEnqueuerQueueIndexFuncCall is an object that describes an invocation
// of method QueueIndex on an instance of MockIndexEnqueuer.
type IndexEnqueuerQueueIndexFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c IndexEnqueuerQueueIndexFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c IndexEnqueuerQueueIndexFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
