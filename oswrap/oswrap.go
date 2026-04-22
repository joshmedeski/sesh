package oswrap

import (
	"os"
)

type Os interface {
	UserConfigDir() (string, error)
	UserHomeDir() (string, error)
	ReadFile(name string) ([]byte, error)
	Getenv(key string) string
	ExpandEnv(s string) string
	Stat(name string) (os.FileInfo, error)
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

func (o *RealOs) Getenv(key string) string {
	return os.Getenv(key)
}

func (o *RealOs) ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func (o *RealOs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
