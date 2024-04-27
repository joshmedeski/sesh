package oswrap

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type Os interface {
	UserHomeDir() (string, error)
}

type RealOs struct{}

func NewOs() Os {
	return &RealOs{}
}

func (o *RealOs) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

type MockOs struct {
	mock.Mock
}

func (m *MockOs) UserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
