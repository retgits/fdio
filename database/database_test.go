// Package database manages storage for fdio
package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	newDBname       = "../test/test.dbx"
	populatedDBname = "../test/test-populated.dbx"
	queries         = []string{"select * from acts where author = \"retgits\"", "select ref, count(*) from acts where author=\"retgits\""}
	statsQueries    = []string{"select author, count(author) as num from acts group by author order by num desc limit 5", "select type, count(type) as num from acts group by type"}
)

func TestDBNew(t *testing.T) {
	fmt.Println("TestDBNew")
	assert := assert.New(t)

	// Try to access a non-existing database file
	db, err := New(newDBname, false, false)
	assert.NotNil(err)
	assert.Equal(err.Error(), fmt.Sprintf("file %s does not exist", newDBname))

	// Try to reset a non-existing database
	db, err = New(newDBname, false, true)
	assert.NotNil(err)
	assert.Equal(err.Error(), fmt.Sprintf("file %s does not exist", newDBname))

	// Create a new database
	db, err = New(newDBname, true, false)
	assert.Equal(newDBname, db.File)
	assert.Nil(err)
}

func TestDBDoQuery(t *testing.T) {
	fmt.Println("TestDBDoQuery")
	assert := assert.New(t)

	db, err := New(populatedDBname, false, false)
	assert.Equal(populatedDBname, db.File)
	assert.Nil(err)

	for _, query := range queries {
		cols, rows, err := db.DoQuery(query)
		assert.Nil(err)
		assert.True(len(cols) > 1)
		assert.True(len(rows) > 0)
	}
}

func TestDBDoStatsQuery(t *testing.T) {
	fmt.Println("TestDBDoStatsQuery")
	assert := assert.New(t)

	db, err := New(populatedDBname, false, false)
	assert.Equal(populatedDBname, db.File)
	assert.Nil(err)

	for _, query := range statsQueries {
		rows, err := db.DoStatsQuery(query)
		assert.Nil(err)
		assert.True(len(rows) > 0)
	}
}

func cleanup() {
	os.Remove(newDBname)
}
