package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbPath = "file::memory:?cache=shared"

func TestCreateEntry(t *testing.T) {
	database := NewSqliteDatabase(dbPath)

	t.Run("Succesful", func(t *testing.T) {
		entry := &Entry{
			Name: "testEntry",
			Path: "/tmp/sesh/",
		}
		err := database.CreateEntry(entry)
		assert.NoError(t, err)
	})

	t.Run("Empty Entry", func(t *testing.T) {
		entry := &Entry{}
		err := database.CreateEntry(entry)
		assert.Error(t, err)
	})

	t.Run("Empty Path", func(t *testing.T) {
		entry := &Entry{
			Name: "empty path",
			Path: "",
		}
		err := database.CreateEntry(entry)
		assert.Equal(t, err, ErrorEntryEmptyPath)
	})

	t.Run("Empty Name", func(t *testing.T) {
		entry := &Entry{
			Name: "",
			Path: "/some/path/",
		}
		err := database.CreateEntry(entry)
		assert.Equal(t, err, ErrorEntryEmptyName)
	})

	t.Run("No Unique Path", func(t *testing.T) {
		entry := &Entry{
			Name: "no unique path",
			Path: "/tmp/sesh/",
		}
		err := database.CreateEntry(entry)
		assert.Error(t, err)
	})

	t.Run("No Unique Name", func(t *testing.T) {
		entry := &Entry{
			Name: "testEntry",
			Path: "/some/random/path/",
		}
		err := database.CreateEntry(entry)
		assert.Error(t, err)
	})
}

func TestUpdateEntry(t *testing.T) {
	database := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "hfpwadsf",
		Path: "new/path/update",
	}
	err := database.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		err := database.UpdateEntry(entry.Name, "name", "updated")
		assert.NoError(t, err)
	})

	t.Run("Update ID", func(t *testing.T) {
		err := database.UpdateEntry(entry.Name, "id", "23468")
		assert.Equal(t, err, ErrorEntryCantUpdateID)
	})

	t.Run("Empty Name", func(t *testing.T) {
		err := database.UpdateEntry(entry.Name, "name", "")
		assert.Equal(t, err, ErrorEntryEmptyName)
	})

	t.Run("Empty Path", func(t *testing.T) {
		err := database.UpdateEntry(entry.Name, "path", "")
		assert.Equal(t, err, ErrorEntryEmptyPath)
	})

}

func TestDeleteEntry(t *testing.T) {
	database := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "aslasd",
		Path: "/asdli/asdf/xc",
	}
	err := database.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		err := database.DeleteEntry(entry.Name)
		assert.NoError(t, err)
	})

}

func TestGetByName(t *testing.T) {
	database := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "asldjhfg",
		Path: "/mps/",
	}
	err := database.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		e, err := database.GetByName(entry.Name)
		assert.NoError(t, err)
		assert.Equal(t, entry.Name, e.Name)
	})

	t.Run("Invalid Name", func(t *testing.T) {
		_, err := database.GetByName("q38fhaalsdf")
		assert.Error(t, err)
	})
}

func TestGetAllEntries(t *testing.T) {
	database := NewSqliteDatabase(dbPath)

	entries, err := database.GetAllEntries()
	assert.NoError(t, err)
	assert.NotEqual(t, []Entry{}, entries)

}
