// Package database manages storage
package database

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/olekukonko/tablewriter"

	// The database is sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// Database implements methods to perform operations on the database.
type Database struct {
	File string
	DB   *sqlx.DB
}

// QueryOptions represents the options you can have for a query and how the result will be rendered
type QueryOptions struct {
	// Writer is where the result is sent to
	Writer io.Writer

	// Query is executed on the database
	Query string

	// MergeCells enables the merge of cells with identical values
	MergeCells bool

	// RowLine enables a line on each row of the table
	RowLine bool

	// Caption sets a caption for the table
	Caption string

	// Render enables the rendering of the table output
	Render bool
}

// QueryResponse represents the response from a query
type QueryResponse struct {
	Rows        [][]string
	ColumnNames []string
	Table       *tablewriter.Table
}

// Contributions is a slice of contribution objects
type Contributions []Contribution

// Contribution is an activity or trigger created for Flogo
type Contribution struct {
	Ref              string `json:"ref"`
	Name             string `json:"name"`
	ContributionType string `json:"type"`
	SourceURL        string
	Author           string `json:"author"`
	UploadedOn       string
	ShowcaseEnabled  bool
	Description      string `json:"description"`
	Version          string `json:"version"`
	Title            string `json:"title"`
	Homepage         string `json:"homepage"`
	Legacy           bool
}

// OpenSession creates a new reference to an SQLite database. If the file cannot be found an exception will be returned.
func OpenSession(file string) (*Database, error) {
	// Validate the file exists
	_, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("error locating database file: %s", err.Error())
	}

	// Connect to the database
	dbase, err := sqlx.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %s", err.Error())
	}

	return &Database{File: file, DB: dbase}, nil
}

// MustOpenSession is like OpenSession but panics if the session cannot be created.
func MustOpenSession(file string) *Database {
	db, err := OpenSession(file)
	if err != nil {
		panic(err)
	}
	return db
}

// Initialize creates the new database structure. This method must be called if you're starting with a brand new database.
func (db *Database) Initialize() error {
	return db.Exec(`create table contributions(
		ref text, 
		name text, 
		contributiontype text, 
		sourceurl text, 
		author text, 
		uploadedon text, 
		showcaseenabled text, 
		description text, 
		version text, 
		title text, 
		homepage text, 
		legacy text)
	`)
}

// Close closes the database and prevents new queries from starting. Close then waits for all queries that have started processing on the server to finish.
func (db *Database) Close() error {
	return db.DB.Close()
}

// Exec executes a query without returning any rows. An error is returned only when the database throws an error.
func (db *Database) Exec(query string) error {
	_, err := db.DB.Exec(query)
	return err
}

// InsertContribution inserts activities and triggers into the database,
func (db *Database) InsertContribution(c Contribution) error {
	q := fmt.Sprintf("insert into contributions(ref, name, contributiontype, sourceurl, author, uploadedon, showcaseenabled, description, version, title, homepage, legacy) values(\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", c.Ref, c.Name, c.ContributionType, c.SourceURL, c.Author, c.UploadedOn, strconv.FormatBool(c.ShowcaseEnabled), c.Description, c.Version, c.Title, c.Homepage, strconv.FormatBool(c.Legacy))
	return db.Exec(q)
}

// Query run a query on the database and prints the result in a table.
func (db *Database) Query(opts QueryOptions) (QueryResponse, error) {
	queryResponse := QueryResponse{}

	// Open a connection to the database
	//dbase, err := sqlx.Open("sqlite3", db.File)
	//if err != nil {
	//	return queryResponse, fmt.Errorf("error while opening connection to database: %s", err.Error())
	//}
	//defer dbase.Close()

	// Execute the query
	rows, err := db.DB.Queryx(opts.Query)
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
