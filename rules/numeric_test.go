package rules

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestRules_max(t *testing.T) {
	max, ok := GetRule("max")

	assert.True(t, ok)
	err := max(&ValidationOptions{
		Value:     1,
		Arguments: []string{"0"},
		Request:   nil,
		Name:      "",
	})

	spew.Dump(err.Error())
	os.Exit(1)
	assert.NoError(t, err)
}
