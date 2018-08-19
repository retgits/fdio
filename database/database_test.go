// Package database manages storage for fdio
package database_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/retgits/fdio/database"
	"github.com/stretchr/testify/assert"
)

var (
	newDBname       = "../test/test.db"
	populatedDBname = "../test/populated.db"
	queries         = []string{"select * from acts where author = \"retgits\"", "select ref, count(*) from acts where author=\"retgits\"", "select author, count(author) as num from acts group by author order by num desc limit 5", "select type, count(type) as num from acts group by type"}
)

func TestDBNew(t *testing.T) {
	fmt.Println("TestDBNew")
	assert := assert.New(t)

	// Try to access a non-existing database file
	db, err := database.New(newDBname, false)
	assert.NotNil(err)
	assert.Equal(err.Error(), "error while checking for database file: stat ../test/test.db: no such file or directory")

	// Create a new database
	db, err = database.New(newDBname, true)
	assert.Equal(newDBname, db.File)
	assert.Nil(err)
	os.Remove(newDBname)
}

func TestDBDoQuery(t *testing.T) {
	fmt.Println("TestDBDoQuery")
	assert := assert.New(t)

	db, err := database.New(populatedDBname, false)
	assert.Equal(populatedDBname, db.File)
	assert.Nil(err)

	for _, query := range queries {
		queryOpts := database.QueryOptions{
			Writer:     os.Stdout,
			Query:      query,
			MergeCells: true,
			RowLine:    true,
			Render:     true,
		}
		response, err := db.RunQuery(queryOpts)
		assert.Nil(err)
		assert.True(len(response.ColumnNames) > 1)
		assert.True(len(response.Rows) > 0)
	}
}

func cleanup() {
	os.Remove(newDBname)
}
