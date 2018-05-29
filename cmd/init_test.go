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
	goCommand   = "go"
	initCommand = []string{"run", "../main.go", "init"}
)

func TestInit(t *testing.T) {
	fmt.Println("TestInit")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// no flags set
	currentCmd := initCommand
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "required flag(s) \"db\" not set")
	outbuf.Reset()
	errbuf.Reset()
}

func TestInitDB(t *testing.T) {
	fmt.Println("TestInitDB")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db flags set, but not create
	currentCmd := append(initCommand, "--db", "./test/test.db")
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stderr := errbuf.String()
	assert.Contains(stderr, "./test/test.db does not exist")
	outbuf.Reset()
	errbuf.Reset()
}

func TestInitDBCreate(t *testing.T) {
	fmt.Println("TestInitDBCreate")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db and create flags set
	currentCmd := append(initCommand, "--db", "../test/test.db", "--create")
	cmd := exec.Command(goCommand, currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	file, err := os.Stat("../test/test.db")
	assert.Equal(file.Name(), "test.db")
	assert.False(os.IsNotExist(err))
	outbuf.Reset()
	errbuf.Reset()
	os.Remove("../test/test.db")
}
