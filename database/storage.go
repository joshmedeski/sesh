package db

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrorEntryCantUpdateID = errors.New("You Cant't Update The ID")
	ErrorEntryEmptyName    = errors.New("No Name Provided")
	ErrorEntryEmptyPath    = errors.New("No Path Provided")
	ErrorEntryNotFound     = gorm.ErrRecordNotFound
)

type Entry struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;index"`
	Path     string `gorm:"unique;index"`
	Score    uint   `gorm:"default:0"`
	Metadata string
}

type Storage interface {
	CreateEntry(*Entry) error
	UpdateEntry(*Entry, string, interface{}) error
	DeleteEntry(uint) error
	GetByID(uint) (Entry, error)
	GetByName(string) (Entry, error)
	GetByMetadata(string) ([]Entry, error)
}
