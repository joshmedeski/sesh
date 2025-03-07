// Code generated by mockery v2.53.0. DO NOT EDIT.

package json

import (
	model "github.com/joshmedeski/sesh/v2/model"
	mock "github.com/stretchr/testify/mock"
)

// MockJson is an autogenerated mock type for the Json type
type MockJson struct {
	mock.Mock
}

type MockJson_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJson) EXPECT() *MockJson_Expecter {
	return &MockJson_Expecter{mock: &_m.Mock}
}

// EncodeSessions provides a mock function with given fields: sessions
func (_m *MockJson) EncodeSessions(sessions []model.SeshSession) string {
	ret := _m.Called(sessions)

	if len(ret) == 0 {
		panic("no return value specified for EncodeSessions")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func([]model.SeshSession) string); ok {
		r0 = rf(sessions)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockJson_EncodeSessions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EncodeSessions'
type MockJson_EncodeSessions_Call struct {
	*mock.Call
}

// EncodeSessions is a helper method to define mock.On call
//   - sessions []model.SeshSession
func (_e *MockJson_Expecter) EncodeSessions(sessions interface{}) *MockJson_EncodeSessions_Call {
	return &MockJson_EncodeSessions_Call{Call: _e.mock.On("EncodeSessions", sessions)}
}

func (_c *MockJson_EncodeSessions_Call) Run(run func(sessions []model.SeshSession)) *MockJson_EncodeSessions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]model.SeshSession))
	})
	return _c
}

func (_c *MockJson_EncodeSessions_Call) Return(_a0 string) *MockJson_EncodeSessions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockJson_EncodeSessions_Call) RunAndReturn(run func([]model.SeshSession) string) *MockJson_EncodeSessions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockJson creates a new instance of MockJson. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJson(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJson {
	mock := &MockJson{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
