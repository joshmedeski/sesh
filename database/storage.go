package db

import (
	"errors"
)

var (
	ErrorEntryCantUpdateID = errors.New("You Cant't Update The ID")
	ErrorEntryEmptyName    = errors.New("No Name Provided")
	ErrorEntryEmptyPath    = errors.New("No Path Provided")
)

type Entry struct {
	ID    uint   `db:"id"`
	Name  string `db:"name"`
	Path  string `db:"path"`
	Score uint   `db:"score"`
}

type Storage interface {
	CreateEntry(*Entry) error
	UpdateEntry(string, string, interface{}) error
	DeleteEntry(string) error
	GetAllEntries() ([]Entry, error)
	GetByName(string) (Entry, error)
}
