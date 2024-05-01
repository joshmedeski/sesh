package runtimewrap

import (
	"runtime"

	"github.com/stretchr/testify/mock"
)

type Runtime interface {
	GOOS() string
}

type RealRunTime struct{}

func NewRunTime() Runtime {
	return &RealRunTime{}
}

func (r *RealRunTime) GOOS() string {
	return runtime.GOOS
}

type MockRunTime struct {
	mock.Mock
}

func (m *MockRunTime) GOOS() string {
	args := m.Called()
	return args.String(0)
}
