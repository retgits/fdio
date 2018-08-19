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
	mainCommand = []string{"run", "../main.go"}
)

func TestMain(t *testing.T) {
	fmt.Println("TestMain")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// no flags set
	currentCmd := mainCommand
	cmd := exec.Command("go", currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "A command-line interface for the Flogo Dot IO website")
	outbuf.Reset()
	errbuf.Reset()
}

func TestVersion(t *testing.T) {
	fmt.Println("TestVersion")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db flags set, but not toml
	currentCmd := append(mainCommand, "--version")
	cmd := exec.Command("go", currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	assert.Contains(stdout, "You're running FDIO version 0.1.0")
	outbuf.Reset()
	errbuf.Reset()
}
