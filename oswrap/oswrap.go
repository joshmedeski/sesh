package oswrap

import (
	"os"
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
