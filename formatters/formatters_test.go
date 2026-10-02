package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		wantType Formatter
	}{
		{name: "stylish", format: "stylish", wantType: &Stylish{}},
		{name: "plain", format: "plain", wantType: &Plain{}},
		{name: "json", format: "json", wantType: &JSON{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.format)
			require.NoError(t, err)
			assert.IsType(t, tt.wantType, got)
		})
	}
}

func TestNewUnknownFormat(t *testing.T) {
	got, err := New("unknown")

	require.Error(t, err)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, `unknown format: "unknown"`)
	assert.ErrorContains(t, err, "available: json, plain, stylish")
}

func TestNewIsCaseSensitive(t *testing.T) {
	_, err := New("Plain")

	assert.Error(t, err)
}

func TestNames(t *testing.T) {
	assert.Equal(t, []string{"json", "plain", "stylish"}, Names())
}
