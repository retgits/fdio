// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FDIOCommandsTestSuite struct {
	suite.Suite
	Command []string
}

func (suite *FDIOCommandsTestSuite) SetupTest() {
	suite.Command = []string{"go", "run", "../main.go"}
}

func (suite *FDIOCommandsTestSuite) TearDownTest() {
	os.Remove("./init.db")
}

func (suite *FDIOCommandsTestSuite) TestRunMain() {
	res, err := runner(suite.Command)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), res, "A command-line interface for the Flogo Dot IO website")
}

func (suite *FDIOCommandsTestSuite) TestRunVersion() {
	args := append(suite.Command, "--version")
	res, err := runner(args)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), res, fmt.Sprintf("You're running FDIO version %s", Version))
}

func (suite *FDIOCommandsTestSuite) TestRunInit() {
	args := append(suite.Command, "init")
	res, err := runner(args)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), res, "Error: required flag(s) \"db\"")

	args = append(args, "--db", "./init.db")
	res, err = runner(args)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), res, "error locating database file: stat ./init.db: no such file or directory")

	os.Create("./init.db")
	res, err = runner(args)
	assert.NoError(suite.T(), err)
}

func (suite *FDIOCommandsTestSuite) TestRunStats() {
	args := append(suite.Command, "stats")
	res, err := runner(args)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), res, "Error: required flag(s) \"db\"")

	args = append(args, "--db", "../test/populated.dbtest")
	res, err = runner(args)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), res, "retgits |   1")
}

func (suite *FDIOCommandsTestSuite) TestRunQuery() {
	args := append(suite.Command, "query")
	res, err := runner(args)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), res, "Error: required flag(s) \"db\"")

	args = append(args, "--db", "../test/populated.dbtest")
	res, err = runner(args)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), res, "Error: required flag(s) \"query\"")

	args = append(args, "--db", "../test/populated.dbtest", "--query", "select author, count(author) as num from contributions group by author order by num desc limit 5")
	res, err = runner(args)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), res, "retgits |   1")
}

func TestCommands(t *testing.T) {
	suite.Run(t, new(FDIOCommandsTestSuite))
}

func runner(args []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	res, err := cmd.CombinedOutput()
	return string(res), err
}
