package code

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGendiff(t *testing.T) {
	got := GenDiff("./testdata/fixture/before.json", "./testdata/fixture/after.json", "json")
	data, err := ReadFile("./testdata/fixture/resultJson.txt")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	assert.Equal(t, got, strings.TrimSpace(string(data)))
}
