package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintUsage(t *testing.T) {
	err := run([]string{"help"})
	assert.NoError(t, err)

	err = run([]string{"--help"})
	assert.NoError(t, err)

	err = run([]string{"-h"})
	assert.NoError(t, err)
}

func TestRunVersion(t *testing.T) {
	err := run([]string{"version"})
	assert.NoError(t, err)

	err = run([]string{"-version"})
	assert.NoError(t, err)
}

func TestUnknownCommand(t *testing.T) {
	err := run([]string{"nonexistent-command-xyz"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}
