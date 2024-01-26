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
