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
	statsCommand = []string{"run", "../main.go", "stats"}
)

func TestStats(t *testing.T) {
	fmt.Println("TestStats")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// no flags set
	currentCmd := statsCommand
	cmd := exec.Command("go", currentCmd...)
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

func TestStatsDB(t *testing.T) {
	fmt.Println("TestStatsDB")
	assert := assert.New(t)

	var outbuf, errbuf bytes.Buffer

	// db flags set
	currentCmd := append(statsCommand, "--db", "../test/populated.db")
	cmd := exec.Command("go", currentCmd...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		fmt.Println(err.Error())
	}
	stdout := outbuf.String()
	stderr := errbuf.String()
	assert.Contains(stdout, "| retgits |  22 |")
	assert.Contains(stdout, "| activity |  21 |")
	assert.True(len(stderr) == 0)
	outbuf.Reset()
	errbuf.Reset()
}
