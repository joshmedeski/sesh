package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbPath = "file::memory:?cache=shared"

func TestCreateEntry(t *testing.T) {
	db := NewSqliteDatabase(dbPath)

	t.Run("Succesful", func(t *testing.T) {
		entry := &Entry{
			Name:     "testEntry",
			Path:     "/tmp/sesh/",
			Score:    1,
			Metadata: "test,entry,metadata",
		}
		err := db.CreateEntry(entry)
		assert.NoError(t, err)
	})

	t.Run("Empty Entry", func(t *testing.T) {
		entry := &Entry{}
		err := db.CreateEntry(entry)
		assert.Error(t, err)
	})

	t.Run("Empty Path", func(t *testing.T) {
		entry := &Entry{
			Name: "empty path",
			Path: "",
		}
		err := db.CreateEntry(entry)
		assert.Equal(t, err, ErrorEntryEmptyPath)
	})

	t.Run("Empty Name", func(t *testing.T) {
		entry := &Entry{
			Name: "",
			Path: "/some/path/",
		}
		err := db.CreateEntry(entry)
		assert.Equal(t, err, ErrorEntryEmptyName)
	})

	t.Run("No Unique Path", func(t *testing.T) {
		entry := &Entry{
			Name: "no unique path",
			Path: "/tmp/sesh/",
		}
		err := db.CreateEntry(entry)
		assert.Error(t, err)
	})

	t.Run("No Unique Name", func(t *testing.T) {
		entry := &Entry{
			Name: "testEntry",
			Path: "/some/random/path/",
		}
		err := db.CreateEntry(entry)
		assert.Error(t, err)
	})
}

func TestUpdateEntry(t *testing.T) {
	db := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "hfpwadsf",
		Path: "new/path/update",
	}
	err := db.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		err := db.UpdateEntry(entry, "name", "updated")
		assert.NoError(t, err)
	})

	t.Run("Update ID", func(t *testing.T) {
		err := db.UpdateEntry(entry, "id", "23468")
		assert.Equal(t, err, ErrorEntryCantUpdateID)
	})

	t.Run("Empty Name", func(t *testing.T) {
		err := db.UpdateEntry(entry, "name", "")
		assert.Equal(t, err, ErrorEntryEmptyName)
	})

	t.Run("Empty Path", func(t *testing.T) {
		err := db.UpdateEntry(entry, "path", "")
		assert.Equal(t, err, ErrorEntryEmptyPath)
	})

}

func TestDeleteEntry(t *testing.T) {
	db := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "aslasd",
		Path: "/asdli/asdf/xc",
	}
	err := db.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		err := db.DeleteEntry(entry.ID)
		assert.NoError(t, err)
	})

}

func TestGetByID(t *testing.T) {
	db := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "testnry",
		Path: "/mpssh/",
	}
	err := db.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		e, err := db.GetByID(entry.ID)
		assert.NoError(t, err)
		assert.Equal(t, entry.ID, e.ID)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		_, err := db.GetByID(3748523245)
		assert.Error(t, err)
	})
}

func TestGetByName(t *testing.T) {
	db := NewSqliteDatabase(dbPath)
	entry := &Entry{
		Name: "asldjhfg",
		Path: "/mps/",
	}
	err := db.CreateEntry(entry)
	assert.NoError(t, err)

	t.Run("Succesful", func(t *testing.T) {
		e, err := db.GetByName(entry.Name)
		assert.NoError(t, err)
		assert.Equal(t, entry.Name, e.Name)
	})

	t.Run("Invalid Name", func(t *testing.T) {
		_, err := db.GetByName("q38fhaalsdf")
		assert.Error(t, err)
	})
}
