// Code generated by mockery v2.43.0. DO NOT EDIT.

package config

import (
	model "github.com/joshmedeski/sesh/model"
	mock "github.com/stretchr/testify/mock"
)

// MockConfig is an autogenerated mock type for the Config type
type MockConfig struct {
	mock.Mock
}

type MockConfig_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConfig) EXPECT() *MockConfig_Expecter {
	return &MockConfig_Expecter{mock: &_m.Mock}
}

// GetConfig provides a mock function with given fields:
func (_m *MockConfig) GetConfig() (model.Config, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfig")
	}

	var r0 model.Config
	var r1 error
	if rf, ok := ret.Get(0).(func() (model.Config, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() model.Config); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Config)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConfig_GetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfig'
type MockConfig_GetConfig_Call struct {
	*mock.Call
}

// GetConfig is a helper method to define mock.On call
func (_e *MockConfig_Expecter) GetConfig() *MockConfig_GetConfig_Call {
	return &MockConfig_GetConfig_Call{Call: _e.mock.On("GetConfig")}
}

func (_c *MockConfig_GetConfig_Call) Run(run func()) *MockConfig_GetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_GetConfig_Call) Return(_a0 model.Config, _a1 error) *MockConfig_GetConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConfig_GetConfig_Call) RunAndReturn(run func() (model.Config, error)) *MockConfig_GetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConfig creates a new instance of MockConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConfig {
	mock := &MockConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}