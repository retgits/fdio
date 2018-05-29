// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	importCommand = []string{"run", "../main.go", "import"}
)

func TestImport(t *testing.T) {
	fmt.Println("TestImport")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// no flags set
	currentCmd := importCommand
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "required flag(s) \"db\", \"toml\" not set")
	outbuf.Reset()
	errbuf.Reset()
}

func TestImportDB(t *testing.T) {
	fmt.Println("TestImportDB")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db flags set, but not toml
	currentCmd := append(importCommand, "--db", "../test/test-import.db")
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "required flag(s) \"toml\" not set")
	outbuf.Reset()
	errbuf.Reset()
}

func TestImportDBWithToml(t *testing.T) {
	fmt.Println("TestImportDBWithToml")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// Create a database to import into
	cmdDB := exec.Command(goCommand, "run", "../main.go", "init", "--db", "../test/test-import.db", "--create")
	errDB := cmdDB.Run()
	if errDB != nil && !strings.Contains(errDB.Error(), "exit status 1") {
		fmt.Println(errDB.Error())
	}

	// db and create flags set
	currentCmd := append(importCommand, "--db", "../test/test-import.db", "--toml", "../test/items.toml")
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	file, err := os.Stat("../test/test-import.db")
	assert.Equal(file.Name(), "test-import.db")
	assert.False(os.IsNotExist(err))
	assert.True(file.Size() > 28000)
	assert.True(file.Size() < 29000)
	outbuf.Reset()
	errbuf.Reset()
	os.Remove("../test/test-import.db")
}
