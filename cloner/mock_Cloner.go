// Code generated by mockery v2.45.1. DO NOT EDIT.

package cloner

import mock "github.com/stretchr/testify/mock"

// MockCloner is an autogenerated mock type for the Cloner type
type MockCloner struct {
	mock.Mock
}

type MockCloner_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCloner) EXPECT() *MockCloner_Expecter {
	return &MockCloner_Expecter{mock: &_m.Mock}
}

// Clone provides a mock function with given fields: path
func (_m *MockCloner) Clone(path string) (string, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for Clone")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCloner_Clone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clone'
type MockCloner_Clone_Call struct {
	*mock.Call
}

// Clone is a helper method to define mock.On call
//   - path string
func (_e *MockCloner_Expecter) Clone(path interface{}) *MockCloner_Clone_Call {
	return &MockCloner_Clone_Call{Call: _e.mock.On("Clone", path)}
}

func (_c *MockCloner_Clone_Call) Run(run func(path string)) *MockCloner_Clone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCloner_Clone_Call) Return(_a0 string, _a1 error) *MockCloner_Clone_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCloner_Clone_Call) RunAndReturn(run func(string) (string, error)) *MockCloner_Clone_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCloner creates a new instance of MockCloner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCloner(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCloner {
	mock := &MockCloner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
