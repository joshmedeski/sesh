package db

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqliteDatabase struct {
	*gorm.DB
}

func NewSqliteDatabase(path string) *SqliteDatabase {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// TODO handle database error
		panic(err)
	}
	db.AutoMigrate(&Entry{})
	return &SqliteDatabase{db}
}

func (db *SqliteDatabase) CreateEntry(entry *Entry) error {
	if entry.Name == "" {
		return ErrorEntryEmptyName
	}
	if entry.Path == "" {
		return ErrorEntryEmptyPath
	}
	return db.Create(entry).Error
}

func (db *SqliteDatabase) UpdateEntry(entry *Entry, column string, value interface{}) error {
	if column == "id" {
		return ErrorEntryCantUpdateID
	}
	if column == "name" && value == "" {
		return ErrorEntryEmptyName
	}
	if column == "path" && value == "" {
		return ErrorEntryEmptyPath
	}
	return db.Model(entry).Where("id = ?", entry.ID).Update(column, value).Error
}

func (db *SqliteDatabase) DeleteEntry(id uint) error {
	return db.Delete(&Entry{}, id).Error
}

func (db *SqliteDatabase) GetByID(id uint) (Entry, error) {
	var entry Entry
	res := db.First(&entry, id)
	if errors.Is(res.Error, ErrorEntryNotFound) {
		return Entry{}, ErrorEntryNotFound
	}
	return entry, nil
}

func (db *SqliteDatabase) GetByName(name string) (Entry, error) {
	var entry Entry
	res := db.First(&entry, "name = ?", name)
	if errors.Is(res.Error, ErrorEntryNotFound) {
		return Entry{}, ErrorEntryNotFound
	}
	return entry, nil
}

// TODO Implement GetByMetadata
func (db *SqliteDatabase) GetByMetadata([]string) ([]Entry, error) { return []Entry{}, nil }
