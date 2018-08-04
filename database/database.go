// Package database manages storage for fdio
package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/olekukonko/tablewriter"

	// The database is sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database and implements methods to perform operations on the database.
type Database struct {
	File string
}

// QueryOptions represents the options you can have for a query and how the result will be rendered
type QueryOptions struct {
	Writer     io.Writer
	Query      string
	MergeCells bool
	RowLine    bool
	Caption    string
	Render     bool
}

// QueryResponse represents the response from a query
type QueryResponse struct {
	Rows        [][]string
	ColumnNames []string
	Table       *tablewriter.Table
}

// New creates a connection to the database. filename represents the exact file location of the database
// file, create a boolean to indiciate whether to create a new file or not if the filename doesn't exist,
// and reset a boolean that indicates whether to delete the existing file and create a new one.
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

// reinit creates the table structure needed in the database. This method must be called if you're
// starting with a brand new database.
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

// InsertContributions inserts activities and triggers into the database. The input argument is an array of
// map[string]interface{} which will be used in the insert statement. Inserts are done in a transaction.
func (db *Database) InsertContributions(items []map[string]interface{}) error {
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
		if item["ref"] != nil {
			// Get values from the item or assign a default value
			ref := getValue(item["ref"], "")
			name := getValue(item["name"], "")
			contribType := getValue(item["type"], "")
			url := getValue(item["url"], "")
			author := getValue(item["author"], "")
			uploadedOn := getValue(item["uploadedon"], "")
			showcase := getValue(item["showcase"], "false")
			description := getValue(item["description"], "")

			// Execute the prepared statement
			_, err = stmt.Exec(ref, name, contribType, description, url, uploadedOn, author, showcase)
			if err != nil {
				// If the ref field already exists in the database, we'll try to update the values assuming the ref field is a valid Go package
				if strings.Contains(err.Error(), "UNIQUE constraint failed: acts.ref") && len(strings.Split(ref, "/")) > 2 {
					urlComponents := strings.Split(ref, "/")
					refToURL := fmt.Sprintf("https://%s/tree/master/%s/", strings.Join(urlComponents[:3], "/"), strings.Join(urlComponents[3:], "/"))

					// Only perform an update if the ref field matches the URL, otherwise it is a fork and the ref should have been updated
					if strings.Contains(refToURL, url) {
						// Create another prepared statement
						updateStmt, err := tx.Prepare("update acts set type=?, description=?, url=?, uploadedon=?, showcase=? where ref=?")
						if err != nil {
							return fmt.Errorf("error while creating update statement: %s", err.Error())
						}
						defer updateStmt.Close()

						// Execute the update statement
						_, err = updateStmt.Exec(contribType, description, url, uploadedOn, showcase, ref)
						if err != nil {
							log.Printf("Error while updating %s: %s\n", ref, err.Error())
						}
					}
				} else {
					log.Printf("Error while inserting %s into database: %s\n", ref, err.Error())
				}
			}
		}
	}

	// Commit the transaction
	tx.Commit()

	return nil
}

func getValue(value interface{}, fallback string) string {
	if value == nil {
		return fallback
	}

	return value.(string)
}

// RunQuery run a query on the database and prints the result in a table
func (db *Database) RunQuery(opts QueryOptions) (QueryResponse, error) {
	queryResponse := QueryResponse{}

	// Open a connection to the database
	dbase, err := sqlx.Open("sqlite3", db.File)
	if err != nil {
		return queryResponse, fmt.Errorf("error while opening connection to database: %s", err.Error())
	}
	defer dbase.Close()

	// Execute the query
	rows, err := dbase.Queryx(opts.Query)
	if err != nil {
		return queryResponse, fmt.Errorf("error while executing query: %s", err.Error())
	}
	defer rows.Close()

	// Get the column names
	colnames, _ := rows.Columns()

	// Prepare the output table
	table := tablewriter.NewWriter(opts.Writer)
	table.SetHeader(colnames)
	table.SetAutoMergeCells(opts.MergeCells)
	table.SetRowLine(opts.RowLine)
	if len(opts.Caption) > 0 {
		table.SetCaption(true, opts.Caption)
	}

	// Prepare a result array
	var resultArray [][]string

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
		table.Append(tempStringArray)
		resultArray = append(resultArray, tempStringArray)
	}

	// Print the table
	if opts.Render {
		table.Render()
	}

	queryResponse.ColumnNames = colnames
	queryResponse.Rows = resultArray
	queryResponse.Table = table

	return queryResponse, nil
}
