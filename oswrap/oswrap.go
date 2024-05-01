package oswrap

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type Os interface {
	UserConfigDir() (string, error)
	UserHomeDir() (string, error)
	ReadFile(name string) ([]byte, error)
}

type RealOs struct{}

func NewOs() Os {
	return &RealOs{}
}

func (o *RealOs) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (o *RealOs) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o *RealOs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

type MockOs struct {
	mock.Mock
}

func (m *MockOs) UserConfigDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOs) UserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOs) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}
