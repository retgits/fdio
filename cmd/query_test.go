// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	queryCommand = []string{"run", "../main.go", "query"}
)

func TestQuery(t *testing.T) {
	fmt.Println("TestQuery")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// no flags set
	currentCmd := queryCommand
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "required flag(s) \"db\", \"query\" not set")
	outbuf.Reset()
	errbuf.Reset()
}

func TestQueryDB(t *testing.T) {
	fmt.Println("TestQueryDB")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db flags set, but not query
	currentCmd := append(queryCommand, "--db", "../test/test-populated.dbx")
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "required flag(s) \"query\" not set")
	outbuf.Reset()
	errbuf.Reset()
}

func TestQueryWithQuery(t *testing.T) {
	fmt.Println("TestQueryWithQuery")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	queries := []string{"select * from acts where author = \"retgits\"", "select ref, count(*) from acts where author=\"retgits\""}

	// db and query flags set
	for _, query := range queries {
		currentCmd := append(queryCommand, "--db", "../test/test-populated.dbx", "--query", query)
		cmd := exec.Command(goCommand, currentCmd...)
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf

		err := cmd.Run()
		if err != nil && !strings.Contains(err.Error(), "exit status 1") {
			fmt.Println(err.Error())
		}
		stdout := outbuf.String()
		assert.Contains(stdout, "retgits/randomstring")
		outbuf.Reset()
		errbuf.Reset()
	}
}
