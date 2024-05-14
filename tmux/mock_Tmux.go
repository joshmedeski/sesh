// Code generated by mockery v2.43.0. DO NOT EDIT.

package tmux

import (
	model "github.com/joshmedeski/sesh/model"
	mock "github.com/stretchr/testify/mock"
)

// MockTmux is an autogenerated mock type for the Tmux type
type MockTmux struct {
	mock.Mock
}

type MockTmux_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTmux) EXPECT() *MockTmux_Expecter {
	return &MockTmux_Expecter{mock: &_m.Mock}
}

// ListSessions provides a mock function with given fields:
func (_m *MockTmux) ListSessions() ([]*model.TmuxSession, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListSessions")
	}

	var r0 []*model.TmuxSession
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*model.TmuxSession, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*model.TmuxSession); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TmuxSession)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTmux_ListSessions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListSessions'
type MockTmux_ListSessions_Call struct {
	*mock.Call
}

// ListSessions is a helper method to define mock.On call
func (_e *MockTmux_Expecter) ListSessions() *MockTmux_ListSessions_Call {
	return &MockTmux_ListSessions_Call{Call: _e.mock.On("ListSessions")}
}

func (_c *MockTmux_ListSessions_Call) Run(run func()) *MockTmux_ListSessions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTmux_ListSessions_Call) Return(_a0 []*model.TmuxSession, _a1 error) *MockTmux_ListSessions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTmux_ListSessions_Call) RunAndReturn(run func() ([]*model.TmuxSession, error)) *MockTmux_ListSessions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTmux creates a new instance of MockTmux. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTmux(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTmux {
	mock := &MockTmux{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}