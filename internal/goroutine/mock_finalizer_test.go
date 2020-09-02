// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.

package goroutine

import "sync"

// MockFinalizer is a mock implementation of the Finalizer interface (from
// the package github.com/sourcegraph/sourcegraph/internal/goroutine) used
// for unit testing.
type MockFinalizer struct {
	// OnShutdownFunc is an instance of a mock function object controlling
	// the behavior of the method OnShutdown.
	OnShutdownFunc *FinalizerOnShutdownFunc
}

// NewMockFinalizer creates a new mock of the Finalizer interface. All
// methods return zero values for all results, unless overwritten.
func NewMockFinalizer() *MockFinalizer {
	return &MockFinalizer{
		OnShutdownFunc: &FinalizerOnShutdownFunc{
			defaultHook: func() {
				return
			},
		},
	}
}

// NewMockFinalizerFrom creates a new mock of the MockFinalizer interface.
// All methods delegate to the given implementation, unless overwritten.
func NewMockFinalizerFrom(i Finalizer) *MockFinalizer {
	return &MockFinalizer{
		OnShutdownFunc: &FinalizerOnShutdownFunc{
			defaultHook: i.OnShutdown,
		},
	}
}

// FinalizerOnShutdownFunc describes the behavior when the OnShutdown method
// of the parent MockFinalizer instance is invoked.
type FinalizerOnShutdownFunc struct {
	defaultHook func()
	hooks       []func()
	history     []FinalizerOnShutdownFuncCall
	mutex       sync.Mutex
}

// OnShutdown delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockFinalizer) OnShutdown() {
	m.OnShutdownFunc.nextHook()()
	m.OnShutdownFunc.appendCall(FinalizerOnShutdownFuncCall{})
	return
}

// SetDefaultHook sets function that is called when the OnShutdown method of
// the parent MockFinalizer instance is invoked and the hook queue is empty.
func (f *FinalizerOnShutdownFunc) SetDefaultHook(hook func()) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// OnShutdown method of the parent MockFinalizer instance inovkes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *FinalizerOnShutdownFunc) PushHook(hook func()) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *FinalizerOnShutdownFunc) SetDefaultReturn() {
	f.SetDefaultHook(func() {
		return
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *FinalizerOnShutdownFunc) PushReturn() {
	f.PushHook(func() {
		return
	})
}

func (f *FinalizerOnShutdownFunc) nextHook() func() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *FinalizerOnShutdownFunc) appendCall(r0 FinalizerOnShutdownFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of FinalizerOnShutdownFuncCall objects
// describing the invocations of this function.
func (f *FinalizerOnShutdownFunc) History() []FinalizerOnShutdownFuncCall {
	f.mutex.Lock()
	history := make([]FinalizerOnShutdownFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// FinalizerOnShutdownFuncCall is an object that describes an invocation of
// method OnShutdown on an instance of MockFinalizer.
type FinalizerOnShutdownFuncCall struct{}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c FinalizerOnShutdownFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c FinalizerOnShutdownFuncCall) Results() []interface{} {
	return []interface{}{}
}
