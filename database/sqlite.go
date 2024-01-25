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
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// TODO handle database error
		panic(err)
	}
	database.AutoMigrate(&Entry{})
	return &SqliteDatabase{database}
}

func (database *SqliteDatabase) CreateEntry(entry *Entry) error {
	if entry.Name == "" {
		return ErrorEntryEmptyName
	}
	if entry.Path == "" {
		return ErrorEntryEmptyPath
	}
	return database.Create(entry).Error
}

func (database *SqliteDatabase) UpdateEntry(entry *Entry, column string, value interface{}) error {
	if column == "id" {
		return ErrorEntryCantUpdateID
	}
	if column == "name" && value == "" {
		return ErrorEntryEmptyName
	}
	if column == "path" && value == "" {
		return ErrorEntryEmptyPath
	}
	return database.Model(entry).Where("id = ?", entry.ID).Update(column, value).Error
}

func (database *SqliteDatabase) DeleteEntry(id uint) error {
	return database.Delete(&Entry{}, id).Error
}

func (database *SqliteDatabase) GetAllEntries() ([]Entry, error) {
	var entries []Entry
	err := database.Find(&entries).Error
	if err != nil {
		return []Entry{}, err
	}
	return entries, nil
}

func (database *SqliteDatabase) GetByID(id uint) (Entry, error) {
	var entry Entry
	res := database.First(&entry, id)
	if errors.Is(res.Error, ErrorEntryNotFound) {
		return Entry{}, ErrorEntryNotFound
	}
	return entry, nil
}

func (database *SqliteDatabase) GetByName(name string) (Entry, error) {
	var entry Entry
	res := database.First(&entry, "name = ?", name)
	if errors.Is(res.Error, ErrorEntryNotFound) {
		return Entry{}, ErrorEntryNotFound
	}
	return entry, nil
}

// TODO Implement GetByMetadata
func (database *SqliteDatabase) GetByMetadata([]string) ([]Entry, error) { return []Entry{}, nil }
