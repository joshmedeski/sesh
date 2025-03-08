// Code generated by mockery v2.52.3. DO NOT EDIT.

package tmuxinator

import (
	model "github.com/joshmedeski/sesh/v2/model"
	mock "github.com/stretchr/testify/mock"
)

// MockTmuxinator is an autogenerated mock type for the Tmuxinator type
type MockTmuxinator struct {
	mock.Mock
}

type MockTmuxinator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTmuxinator) EXPECT() *MockTmuxinator_Expecter {
	return &MockTmuxinator_Expecter{mock: &_m.Mock}
}

// List provides a mock function with no fields
func (_m *MockTmuxinator) List() ([]*model.TmuxinatorConfig, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*model.TmuxinatorConfig
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*model.TmuxinatorConfig, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*model.TmuxinatorConfig); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TmuxinatorConfig)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTmuxinator_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockTmuxinator_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
func (_e *MockTmuxinator_Expecter) List() *MockTmuxinator_List_Call {
	return &MockTmuxinator_List_Call{Call: _e.mock.On("List")}
}

func (_c *MockTmuxinator_List_Call) Run(run func()) *MockTmuxinator_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTmuxinator_List_Call) Return(_a0 []*model.TmuxinatorConfig, _a1 error) *MockTmuxinator_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTmuxinator_List_Call) RunAndReturn(run func() ([]*model.TmuxinatorConfig, error)) *MockTmuxinator_List_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: targetSession
func (_m *MockTmuxinator) Start(targetSession string) (string, error) {
	ret := _m.Called(targetSession)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(targetSession)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(targetSession)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(targetSession)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTmuxinator_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type MockTmuxinator_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - targetSession string
func (_e *MockTmuxinator_Expecter) Start(targetSession interface{}) *MockTmuxinator_Start_Call {
	return &MockTmuxinator_Start_Call{Call: _e.mock.On("Start", targetSession)}
}

func (_c *MockTmuxinator_Start_Call) Run(run func(targetSession string)) *MockTmuxinator_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockTmuxinator_Start_Call) Return(_a0 string, _a1 error) *MockTmuxinator_Start_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTmuxinator_Start_Call) RunAndReturn(run func(string) (string, error)) *MockTmuxinator_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTmuxinator creates a new instance of MockTmuxinator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTmuxinator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTmuxinator {
	mock := &MockTmuxinator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
