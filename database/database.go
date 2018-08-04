// Package database manages storage for fdio
package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	// The database is sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database and implements methods to perform
// operations on the database.
type Database struct {
	File string
}

// New creates a connection to the database. filename represents the exact file location
// of the database file, create a boolean to indiciate whether to create a new file or
// not if the filename doesn't exist, and reset a boolean that indicates whether to delete
// the existing file and create a new one.
func New(filename string, create bool, reset bool) (*Database, error) {
	// Remove database file if requested
	if reset {
		log.Printf("Reset fdio database...\n")
		err := os.Remove(filename)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("error while removing database file: %s", err.Error())
		}
	}

	// Create database file if requestes
	if create {
		log.Printf("Create new fdio database...\n")
		_, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("error while creating database file: %s", err.Error())
		}
	}

	// Making sure the file exists
	_, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("file %s does not exist", filename)
	}

	db := &Database{File: filename}

	// If the database was reset or newly created, recreate the table structure
	if reset || create {
		db.reinit()
	}

	return db, nil
}

// reinit creates the table structure needed in the database. This method must be called
// if you're starting with a brand new database.
func (db *Database) reinit() error {
	// Open a connection to the database
	dbase, err := sqlx.Open("sqlite3", db.File)
	if err != nil {
		return fmt.Errorf("error while opening connection to database: %s", err.Error())
	}
	defer dbase.Close()

	// Create the new table
	_, err = dbase.Exec("create table acts (ref text not null primary key, name text, type text, description text, url text, uploadedon text, author text, showcase text)")
	if err != nil {
		return fmt.Errorf("error while creating table: %s", err.Error())
	}

	return nil
}

// InsertActs inserts activities and triggers into the database. The input argument is an
// array of map[string]interface{} which will be used in the insert statement. Inserts are
// done in a transaction.
func (db *Database) InsertActs(items []map[string]interface{}) error {
	// Open a connection to the database
	dbase, err := sqlx.Open("sqlite3", db.File)
	if err != nil {
		return fmt.Errorf("error while opening connection to database: %s", err.Error())
	}
	defer dbase.Close()

	// Start a transaction to add everything into the database
	tx, err := dbase.Begin()
	if err != nil {
		return fmt.Errorf("error while starting database transaction: %s", err.Error())
	}

	// Create a prepared statement
	stmt, err := tx.Prepare("insert into acts(ref, name, type, description, url, uploadedon, author, showcase) values(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error while creating prepared statement: %s", err.Error())
	}
	defer stmt.Close()

	// Insert items into database
	for _, item := range items {
		if item["name"] != nil {

			key := item["ref"].(string)

			if len(item["uploadedon"].(string)) == 0 {
				item["uploadedon"] = ""
			}

			if item["showcase"] == nil || item["showcase"].(string) != "true" {
				item["showcase"] = "false"
			}

			_, err = stmt.Exec(key, item["name"].(string), item["type"].(string), item["description"].(string), item["url"].(string), item["uploadedon"].(string), item["author"].(string), item["showcase"].(string))
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed: acts.ref") {
					log.Printf("Key %s already exists, trying to update\n", key)
					urlComponents := strings.Split(key, "/")
					// We can only update valid Go packages...
					if len(urlComponents) > 2 {
						refToURL := fmt.Sprintf("https://%s/tree/master/%s/", strings.Join(urlComponents[:3], "/"), strings.Join(urlComponents[3:], "/"))
						// Only perform an update if the ref field matches the URL, otherwise it is a fork and the ref should have been updated
						if strings.Contains(refToURL, item["url"].(string)) {
							updateStmt, err := tx.Prepare("update acts set type=?, description=?, url=?, uploadedon=?, showcase=? where ref=?")
							if err != nil {
								return fmt.Errorf("error while creating update statement: %s", err.Error())
							}
							defer updateStmt.Close()
							_, err = updateStmt.Exec(item["type"].(string), item["description"].(string), item["url"].(string), item["uploadedon"].(string), item["showcase"].(string), key)
							if err != nil {
								log.Printf("Error while updating %s: %s\n", key, err.Error())
							}
						}
					}
				} else {
					log.Printf("Error while inserting %s into database: %s\n", key, err.Error())
				}
			}
		}
	}

	// Commit the transaction
	tx.Commit()

	return nil
}

// DoQuery executes a query on the database and returns the results
// []string are the column headers
// [][]string are the rows with data fields
func (db *Database) DoQuery(query string) ([]string, [][]string, error) {
	// Open a connection to the database
	dbase, err := sqlx.Open("sqlite3", db.File)
	if err != nil {
		return nil, nil, fmt.Errorf("error while opening connection to database: %s", err.Error())
	}
	defer dbase.Close()

	// Execute the query
	rows, err := dbase.Queryx(query)
	if err != nil {
		return nil, nil, fmt.Errorf("error while executing query: %s", err.Error())
	}
	defer rows.Close()

	// Prepare a result array
	var resultArray [][]string

	// Get the column names
	colnames, _ := rows.Columns()

	// Loop over the result
	for rows.Next() {
		cols, _ := rows.SliceScan()
		tempStringArray := make([]string, len(cols))
		for idx := range cols {
			switch v := cols[idx].(type) {
			case int64:
				tempStringArray[idx] = strconv.Itoa(int(v))
			case string:
				tempStringArray[idx] = v
			case nil:
				tempStringArray[idx] = ""
			default:
				tempStringArray[idx] = string(v.([]uint8))
			}
		}
		resultArray = append(resultArray, tempStringArray)
	}

	return colnames, resultArray, nil
}

// DoStatsQuery executes a stats query on the database and returns the results
func (db *Database) DoStatsQuery(query string) ([]string, error) {
	// Open a connection to the database
	dbase, err := sqlx.Open("sqlite3", db.File)
	if err != nil {
		return nil, fmt.Errorf("error while opening connection to database: %s", err.Error())
	}
	defer dbase.Close()

	// Execute the query
	rows, err := dbase.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %s", err.Error())
	}
	defer rows.Close()

	// Prepare a result array
	var resultArray []string

	// Loop over the resultset
	for rows.Next() {
		result := make([]string, 2)
		err = rows.Scan(&result[0], &result[1])
		if err != nil {
			log.Fatal(err)
		}
		resultArray = append(resultArray, fmt.Sprintf("%s (%s)", result[0], result[1]))
	}

	return resultArray, nil
}
