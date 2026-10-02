package formatters

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStylishRender(t *testing.T) {
	tests := []struct {
		name     string
		before   string
		after    string
		wantFile string
	}{
		{
			name:     "flat json",
			before:   "../testdata/fixture/before.json",
			after:    "../testdata/fixture/after.json",
			wantFile: "../testdata/fixture/stylishFlat.txt",
		},
		{
			name:     "nested json",
			before:   "../testdata/fixture/beforeTree.json",
			after:    "../testdata/fixture/afterTree.json",
			wantFile: "../testdata/fixture/stylishTree.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBytes, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			want := strings.TrimSpace(string(wantBytes))
			got := strings.TrimSpace((&Stylish{}).Render(BuildAST(loadFixture(t, tt.before), loadFixture(t, tt.after))))

			assert.Equal(t, want, got)
		})
	}
}

func TestStylishRenderValues(t *testing.T) {
	nodes := []Node{
		{Name: "map", Type: Added, ValueAfter: map[string]any{"key": "value"}},
		{Name: "slice", Type: Added, ValueAfter: []any{1, 2}},
		{Name: "null", Type: Added, ValueAfter: nil},
		{Name: "emptyMap", Type: Added, ValueAfter: map[string]any{}},
		{Name: "number", Type: Added, ValueAfter: 42.0},
		{Name: "deep", Type: Added, ValueAfter: map[string]any{
			"b": map[string]any{"c": "d"},
			"a": 1.0,
		}},
	}

	want := strings.Join([]string{
		"{",
		"  + map: {",
		"        key: value",
		"    }",
		"  + slice: [complex value]",
		"  + null: null",
		"  + emptyMap: {}",
		"  + number: 42",
		"  + deep: {",
		"        a: 1",
		"        b: {",
		"            c: d",
		"        }",
		"    }",
		"}",
	}, "\n")

	assert.Equal(t, want, (&Stylish{}).Render(nodes))
}

func loadFixture(t *testing.T, path string) map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var out map[string]any
	require.NoError(t, json.Unmarshal(data, &out))

	return out
}
