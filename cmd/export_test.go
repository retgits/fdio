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
	exportCommand = []string{"run", "../main.go", "export"}
)

func TestExport(t *testing.T) {
	fmt.Println("TestExport")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// no flags set
	currentCmd := exportCommand
	cmd := exec.Command("go", currentCmd...)
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

func TestExportDB(t *testing.T) {
	fmt.Println("TestExportDB")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db flags set, but not toml
	currentCmd := append(exportCommand, "--db", "../test/test-populated.dbx")
	cmd := exec.Command("go", currentCmd...)
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

func TestExportDBWithToml(t *testing.T) {
	fmt.Println("TestExportDBWithToml")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db and create flags set
	currentCmd := append(exportCommand, "--db", "../test/populated.db", "--toml", "../test/test-export.toml")
	cmd := exec.Command("go", currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	file, err := os.Stat("../test/test-export.toml")
	assert.Equal(file.Name(), "test-export.toml")
	assert.False(os.IsNotExist(err))
	assert.True(file.Size() > 5000)
	assert.True(file.Size() < 10000)
	outbuf.Reset()
	errbuf.Reset()
	os.Remove("../test/test-export.toml")
}
