package db

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDatabase struct {
	*sql.DB
}

func NewSqliteDatabase(path string) *SqliteDatabase {
	database, _ := sql.Open("sqlite3", path)
	database.Exec(`
    CREATE TABLE IF NOT EXISTS entries (
        id INTEGER PRIMARY KEY,
        name TEXT UNIQUE,
        path TEXT UNIQUE,
        score INTEGER DEFAULT 0
    );
    `)
	return &SqliteDatabase{database}
}

func (database *SqliteDatabase) CreateEntry(entry *Entry) error {
	q := "INSERT INTO entries (name, path) VALUES (?, ?);"
	if entry.Name == "" {
		return ErrorEntryEmptyName
	}
	if entry.Path == "" {
		return ErrorEntryEmptyPath
	}
	_, err := database.Exec(q, entry.Name, entry.Path)
	if err != nil {
		return err
	}
	return nil
}

func (database *SqliteDatabase) UpdateEntry(name string, column string, value interface{}) error {
	q := "UPDATE entries SET %column% = ? WHERE name = ?;"
	if column == "id" {
		return ErrorEntryCantUpdateID
	}
	if column == "name" && value == "" {
		return ErrorEntryEmptyName
	}
	if column == "path" && value == "" {
		return ErrorEntryEmptyPath
	}
	_, err := database.Exec(strings.ReplaceAll(q, "%column%", column), value, name)
	if err != nil {
		return err
	}
	return nil
}

func (database *SqliteDatabase) DeleteEntry(name string) error {
	q := "DELETE FROM entries WHERE name = ?;"
	_, err := database.Exec(q, name)
	if err != nil {
		return err
	}
	return nil
}

func (database *SqliteDatabase) GetAllEntries() ([]Entry, error) {
	var entries []Entry
	q := "SELECT * FROM entries;"
	rows, err := database.Query(q)
	if err != nil {
		return []Entry{}, err
	}

	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.Name, &entry.Path, &entry.Score)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (database *SqliteDatabase) GetByName(name string) (Entry, error) {
	var entry Entry
	q := "SELECT * FROM entries WHERE name = ?;"
	res := database.QueryRow(q, name)
	err := res.Scan(&entry.ID, &entry.Name, &entry.Path, &entry.Score)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}
