// Code generated by mockery v2.51.1. DO NOT EDIT.

package previewer

import mock "github.com/stretchr/testify/mock"

// MockPreviewer is an autogenerated mock type for the Previewer type
type MockPreviewer struct {
	mock.Mock
}

type MockPreviewer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPreviewer) EXPECT() *MockPreviewer_Expecter {
	return &MockPreviewer_Expecter{mock: &_m.Mock}
}

// Preview provides a mock function with given fields: name
func (_m *MockPreviewer) Preview(name string) (string, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Preview")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPreviewer_Preview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Preview'
type MockPreviewer_Preview_Call struct {
	*mock.Call
}

// Preview is a helper method to define mock.On call
//   - name string
func (_e *MockPreviewer_Expecter) Preview(name interface{}) *MockPreviewer_Preview_Call {
	return &MockPreviewer_Preview_Call{Call: _e.mock.On("Preview", name)}
}

func (_c *MockPreviewer_Preview_Call) Run(run func(name string)) *MockPreviewer_Preview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPreviewer_Preview_Call) Return(_a0 string, _a1 error) *MockPreviewer_Preview_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPreviewer_Preview_Call) RunAndReturn(run func(string) (string, error)) *MockPreviewer_Preview_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPreviewer creates a new instance of MockPreviewer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPreviewer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPreviewer {
	mock := &MockPreviewer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
