// Code generated by mockery v2.52.3. DO NOT EDIT.

package icon

import (
	model "github.com/joshmedeski/sesh/v2/model"
	mock "github.com/stretchr/testify/mock"
)

// MockIcon is an autogenerated mock type for the Icon type
type MockIcon struct {
	mock.Mock
}

type MockIcon_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIcon) EXPECT() *MockIcon_Expecter {
	return &MockIcon_Expecter{mock: &_m.Mock}
}

// AddIcon provides a mock function with given fields: session
func (_m *MockIcon) AddIcon(session model.SeshSession) string {
	ret := _m.Called(session)

	if len(ret) == 0 {
		panic("no return value specified for AddIcon")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(model.SeshSession) string); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIcon_AddIcon_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddIcon'
type MockIcon_AddIcon_Call struct {
	*mock.Call
}

// AddIcon is a helper method to define mock.On call
//   - session model.SeshSession
func (_e *MockIcon_Expecter) AddIcon(session interface{}) *MockIcon_AddIcon_Call {
	return &MockIcon_AddIcon_Call{Call: _e.mock.On("AddIcon", session)}
}

func (_c *MockIcon_AddIcon_Call) Run(run func(session model.SeshSession)) *MockIcon_AddIcon_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.SeshSession))
	})
	return _c
}

func (_c *MockIcon_AddIcon_Call) Return(_a0 string) *MockIcon_AddIcon_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIcon_AddIcon_Call) RunAndReturn(run func(model.SeshSession) string) *MockIcon_AddIcon_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveIcon provides a mock function with given fields: name
func (_m *MockIcon) RemoveIcon(name string) string {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for RemoveIcon")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIcon_RemoveIcon_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveIcon'
type MockIcon_RemoveIcon_Call struct {
	*mock.Call
}

// RemoveIcon is a helper method to define mock.On call
//   - name string
func (_e *MockIcon_Expecter) RemoveIcon(name interface{}) *MockIcon_RemoveIcon_Call {
	return &MockIcon_RemoveIcon_Call{Call: _e.mock.On("RemoveIcon", name)}
}

func (_c *MockIcon_RemoveIcon_Call) Run(run func(name string)) *MockIcon_RemoveIcon_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIcon_RemoveIcon_Call) Return(_a0 string) *MockIcon_RemoveIcon_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIcon_RemoveIcon_Call) RunAndReturn(run func(string) string) *MockIcon_RemoveIcon_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIcon creates a new instance of MockIcon. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIcon(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIcon {
	mock := &MockIcon{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
