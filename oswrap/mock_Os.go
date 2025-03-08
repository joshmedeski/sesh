// Code generated by mockery v2.52.3. DO NOT EDIT.

package oswrap

import (
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// MockOs is an autogenerated mock type for the Os type
type MockOs struct {
	mock.Mock
}

type MockOs_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOs) EXPECT() *MockOs_Expecter {
	return &MockOs_Expecter{mock: &_m.Mock}
}

// Getenv provides a mock function with given fields: key
func (_m *MockOs) Getenv(key string) string {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Getenv")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockOs_Getenv_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Getenv'
type MockOs_Getenv_Call struct {
	*mock.Call
}

// Getenv is a helper method to define mock.On call
//   - key string
func (_e *MockOs_Expecter) Getenv(key interface{}) *MockOs_Getenv_Call {
	return &MockOs_Getenv_Call{Call: _e.mock.On("Getenv", key)}
}

func (_c *MockOs_Getenv_Call) Run(run func(key string)) *MockOs_Getenv_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockOs_Getenv_Call) Return(_a0 string) *MockOs_Getenv_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOs_Getenv_Call) RunAndReturn(run func(string) string) *MockOs_Getenv_Call {
	_c.Call.Return(run)
	return _c
}

// ReadFile provides a mock function with given fields: name
func (_m *MockOs) ReadFile(name string) ([]byte, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for ReadFile")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOs_ReadFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadFile'
type MockOs_ReadFile_Call struct {
	*mock.Call
}

// ReadFile is a helper method to define mock.On call
//   - name string
func (_e *MockOs_Expecter) ReadFile(name interface{}) *MockOs_ReadFile_Call {
	return &MockOs_ReadFile_Call{Call: _e.mock.On("ReadFile", name)}
}

func (_c *MockOs_ReadFile_Call) Run(run func(name string)) *MockOs_ReadFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockOs_ReadFile_Call) Return(_a0 []byte, _a1 error) *MockOs_ReadFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOs_ReadFile_Call) RunAndReturn(run func(string) ([]byte, error)) *MockOs_ReadFile_Call {
	_c.Call.Return(run)
	return _c
}

// Stat provides a mock function with given fields: name
func (_m *MockOs) Stat(name string) (os.FileInfo, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Stat")
	}

	var r0 os.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (os.FileInfo, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) os.FileInfo); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(os.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOs_Stat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stat'
type MockOs_Stat_Call struct {
	*mock.Call
}

// Stat is a helper method to define mock.On call
//   - name string
func (_e *MockOs_Expecter) Stat(name interface{}) *MockOs_Stat_Call {
	return &MockOs_Stat_Call{Call: _e.mock.On("Stat", name)}
}

func (_c *MockOs_Stat_Call) Run(run func(name string)) *MockOs_Stat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockOs_Stat_Call) Return(_a0 os.FileInfo, _a1 error) *MockOs_Stat_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOs_Stat_Call) RunAndReturn(run func(string) (os.FileInfo, error)) *MockOs_Stat_Call {
	_c.Call.Return(run)
	return _c
}

// UserConfigDir provides a mock function with no fields
func (_m *MockOs) UserConfigDir() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for UserConfigDir")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOs_UserConfigDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserConfigDir'
type MockOs_UserConfigDir_Call struct {
	*mock.Call
}

// UserConfigDir is a helper method to define mock.On call
func (_e *MockOs_Expecter) UserConfigDir() *MockOs_UserConfigDir_Call {
	return &MockOs_UserConfigDir_Call{Call: _e.mock.On("UserConfigDir")}
}

func (_c *MockOs_UserConfigDir_Call) Run(run func()) *MockOs_UserConfigDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockOs_UserConfigDir_Call) Return(_a0 string, _a1 error) *MockOs_UserConfigDir_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOs_UserConfigDir_Call) RunAndReturn(run func() (string, error)) *MockOs_UserConfigDir_Call {
	_c.Call.Return(run)
	return _c
}

// UserHomeDir provides a mock function with no fields
func (_m *MockOs) UserHomeDir() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for UserHomeDir")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOs_UserHomeDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserHomeDir'
type MockOs_UserHomeDir_Call struct {
	*mock.Call
}

// UserHomeDir is a helper method to define mock.On call
func (_e *MockOs_Expecter) UserHomeDir() *MockOs_UserHomeDir_Call {
	return &MockOs_UserHomeDir_Call{Call: _e.mock.On("UserHomeDir")}
}

func (_c *MockOs_UserHomeDir_Call) Run(run func()) *MockOs_UserHomeDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockOs_UserHomeDir_Call) Return(_a0 string, _a1 error) *MockOs_UserHomeDir_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOs_UserHomeDir_Call) RunAndReturn(run func() (string, error)) *MockOs_UserHomeDir_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOs creates a new instance of MockOs. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOs(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOs {
	mock := &MockOs{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
