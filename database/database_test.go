// Package database manages storage for fdio
package database

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DBOpsTestSuite struct {
	suite.Suite
	DatabaseToCreate    string
	NotExistingDatabase string
}

type DBQueryTestSuite struct {
	suite.Suite
	DatabaseToCreate string
	db               *Database
}

func (suite *DBOpsTestSuite) SetupTest() {
	suite.DatabaseToCreate = "./my.db"
	suite.NotExistingDatabase = "./other.db"
	os.Create(suite.DatabaseToCreate)
}

func (suite *DBOpsTestSuite) TearDownTest() {
	os.Remove(suite.DatabaseToCreate)
}

func (suite *DBQueryTestSuite) SetupTest() {
	suite.DatabaseToCreate = "./mycontrib.db"
	os.Create(suite.DatabaseToCreate)
	db, _ := OpenSession(suite.DatabaseToCreate)
	db.Initialize()
	suite.db = db
}

func (suite *DBQueryTestSuite) TearDownTest() {
	os.Remove(suite.DatabaseToCreate)
}

func (suite *DBOpsTestSuite) TestOpenSession() {
	db, err := OpenSession(suite.NotExistingDatabase)
	assert.Nil(suite.T(), db)
	assert.Error(suite.T(), err)

	db, err = OpenSession(suite.DatabaseToCreate)
	assert.NotNil(suite.T(), db)
	assert.NoError(suite.T(), err)

	assert.Panics(suite.T(), func() { MustOpenSession(suite.NotExistingDatabase) })

	db = MustOpenSession(suite.DatabaseToCreate)
	assert.NotNil(suite.T(), db)
}

func (suite *DBOpsTestSuite) TestInitializeDBStructure() {
	db, _ := OpenSession(suite.NotExistingDatabase)

	assert.Panics(suite.T(), func() { db.Initialize() })

	db, _ = OpenSession(suite.DatabaseToCreate)
	err := db.Initialize()
	assert.NoError(suite.T(), err)

	err = db.Initialize()
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, "table contributions already exists")
}

func (suite *DBQueryTestSuite) TestInsertContrib() {
	c := Contribution{
		Author:           "retgits",
		ContributionType: "flogo:activity",
		Description:      "A new awesome contribution",
		Homepage:         "https://flogo.io",
		Name:             "awesomeness",
		Ref:              "deprecated",
		ShowcaseEnabled:  "no",
		SourceURL:        "https://github.com/retgits",
		Title:            "AwesomeContrib",
		UploadedOn:       time.Now(),
		Version:          "0.1.0",
	}
	err := suite.db.InsertContribution(c)
	assert.NoError(suite.T(), err)
}

func (suite *DBQueryTestSuite) TestQuery() {
	c := Contribution{
		Author:           "retgits",
		ContributionType: "flogo:activity",
		Description:      "A new awesome contribution",
		Homepage:         "https://flogo.io",
		Name:             "awesomeness",
		Ref:              "deprecated",
		ShowcaseEnabled:  "no",
		SourceURL:        "https://github.com/retgits",
		Title:            "AwesomeContrib",
		UploadedOn:       time.Now(),
		Version:          "0.1.0",
	}
	suite.db.InsertContribution(c)

	o := QueryOptions{
		Writer:     os.Stdout,
		Caption:    "This table contains all contributions",
		MergeCells: false,
		Query:      "select * from contributions",
		Render:     true,
		RowLine:    true,
	}

	res, err := suite.db.Query(o)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *DBOpsTestSuite) TestCloseDB() {
	db, _ := OpenSession(suite.NotExistingDatabase)

	assert.Panics(suite.T(), func() { db.Close() })

	db, _ = OpenSession(suite.DatabaseToCreate)
	err := db.Close()
	assert.NoError(suite.T(), err)
}

func TestInitTestSuite(t *testing.T) {
	suite.Run(t, new(DBOpsTestSuite))
	suite.Run(t, new(DBQueryTestSuite))
}
