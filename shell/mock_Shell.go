// Code generated by mockery v2.46.0. DO NOT EDIT.

package shell

import mock "github.com/stretchr/testify/mock"

// MockShell is an autogenerated mock type for the Shell type
type MockShell struct {
	mock.Mock
}

type MockShell_Expecter struct {
	mock *mock.Mock
}

func (_m *MockShell) EXPECT() *MockShell_Expecter {
	return &MockShell_Expecter{mock: &_m.Mock}
}

// Cmd provides a mock function with given fields: cmd, arg
func (_m *MockShell) Cmd(cmd string, arg ...string) (string, error) {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, cmd)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Cmd")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...string) (string, error)); ok {
		return rf(cmd, arg...)
	}
	if rf, ok := ret.Get(0).(func(string, ...string) string); ok {
		r0 = rf(cmd, arg...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = rf(cmd, arg...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockShell_Cmd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cmd'
type MockShell_Cmd_Call struct {
	*mock.Call
}

// Cmd is a helper method to define mock.On call
//   - cmd string
//   - arg ...string
func (_e *MockShell_Expecter) Cmd(cmd interface{}, arg ...interface{}) *MockShell_Cmd_Call {
	return &MockShell_Cmd_Call{Call: _e.mock.On("Cmd",
		append([]interface{}{cmd}, arg...)...)}
}

func (_c *MockShell_Cmd_Call) Run(run func(cmd string, arg ...string)) *MockShell_Cmd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockShell_Cmd_Call) Return(_a0 string, _a1 error) *MockShell_Cmd_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockShell_Cmd_Call) RunAndReturn(run func(string, ...string) (string, error)) *MockShell_Cmd_Call {
	_c.Call.Return(run)
	return _c
}

// ListCmd provides a mock function with given fields: cmd, arg
func (_m *MockShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, cmd)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListCmd")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...string) ([]string, error)); ok {
		return rf(cmd, arg...)
	}
	if rf, ok := ret.Get(0).(func(string, ...string) []string); ok {
		r0 = rf(cmd, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = rf(cmd, arg...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockShell_ListCmd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListCmd'
type MockShell_ListCmd_Call struct {
	*mock.Call
}

// ListCmd is a helper method to define mock.On call
//   - cmd string
//   - arg ...string
func (_e *MockShell_Expecter) ListCmd(cmd interface{}, arg ...interface{}) *MockShell_ListCmd_Call {
	return &MockShell_ListCmd_Call{Call: _e.mock.On("ListCmd",
		append([]interface{}{cmd}, arg...)...)}
}

func (_c *MockShell_ListCmd_Call) Run(run func(cmd string, arg ...string)) *MockShell_ListCmd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockShell_ListCmd_Call) Return(_a0 []string, _a1 error) *MockShell_ListCmd_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockShell_ListCmd_Call) RunAndReturn(run func(string, ...string) ([]string, error)) *MockShell_ListCmd_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockShell creates a new instance of MockShell. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockShell(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockShell {
	mock := &MockShell{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
