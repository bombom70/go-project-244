package formatters

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONRender(t *testing.T) {
	tests := []struct {
		name     string
		before   string
		after    string
		wantFile string
	}{
		{
			name:     "nested json files",
			before:   "../testdata/fixture/beforeTree.json",
			after:    "../testdata/fixture/afterTree.json",
			wantFile: "../testdata/fixture/jsonTree.json",
		},
		{
			name:     "flat json files",
			before:   "../testdata/fixture/before.json",
			after:    "../testdata/fixture/after.json",
			wantFile: "../testdata/fixture/jsonFlat.json",
		},
		{
			name:     "official hexlet fixtures",
			before:   "../testdata/fixture/file1.json",
			after:    "../testdata/fixture/file2.json",
			wantFile: "../testdata/fixture/result_json.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBytes, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			want := strings.TrimSpace(string(wantBytes))
			got := (&JSON{}).Render(BuildAST(loadFixture(t, tt.before), loadFixture(t, tt.after)))

			assert.Equal(t, want, got)
			assert.True(t, json.Valid([]byte(got)), "output must be valid json")
		})
	}
}

func TestJSONRenderStatuses(t *testing.T) {
	nodes := []Node{
		{Name: "added", Type: Added, ValueAfter: true},
		{Name: "removed", Type: Deleted, ValueBefore: 200.0},
		{Name: "changed", Type: Changed, ValueBefore: "bas", ValueAfter: "bars"},
		{Name: "same", Type: Unchanged, ValueBefore: "hexlet.io", ValueAfter: "hexlet.io"},
	}

	want := `{
  "key": "",
  "type": "root",
  "children": [
    {
      "key": "added",
      "type": "added",
      "value2": true
    },
    {
      "key": "removed",
      "type": "deleted",
      "value1": 200
    },
    {
      "key": "changed",
      "type": "changed",
      "value1": "bas",
      "value2": "bars"
    },
    {
      "key": "same",
      "type": "unchanged",
      "value1": "hexlet.io"
    }
  ]
}`

	assert.Equal(t, want, (&JSON{}).Render(nodes))
}

func TestJSONRenderNestedAndComplexValues(t *testing.T) {
	nodes := []Node{
		{Name: "common", Type: Nested, Children: []Node{
			{Name: "setting6", Type: Nested, Children: []Node{
				{Name: "doge", Type: Nested, Children: []Node{
					{Name: "wow", Type: Changed, ValueBefore: "", ValueAfter: "so much"},
				}},
				{Name: "key", Type: Unchanged, ValueBefore: "value", ValueAfter: "value"},
			}},
		}},
		{Name: "group2", Type: Deleted, ValueBefore: map[string]any{"abc": 12345.0}},
		{Name: "group3", Type: Added, ValueAfter: map[string]any{"fee": 100500.0}},
		{Name: "hosts", Type: Added, ValueAfter: []any{"a", "b"}},
		{Name: "empty", Type: Nested},
	}

	want := `{
  "key": "",
  "type": "root",
  "children": [
    {
      "key": "common",
      "type": "nested",
      "children": [
        {
          "key": "setting6",
          "type": "nested",
          "children": [
            {
              "key": "doge",
              "type": "nested",
              "children": [
                {
                  "key": "wow",
                  "type": "changed",
                  "value1": "",
                  "value2": "so much"
                }
              ]
            },
            {
              "key": "key",
              "type": "unchanged",
              "value1": "value"
            }
          ]
        }
      ]
    },
    {
      "key": "group2",
      "type": "deleted",
      "value1": {
        "abc": 12345
      }
    },
    {
      "key": "group3",
      "type": "added",
      "value2": {
        "fee": 100500
      }
    },
    {
      "key": "hosts",
      "type": "added",
      "value2": [
        "a",
        "b"
      ]
    },
    {
      "key": "empty",
      "type": "nested",
      "children": []
    }
  ]
}`

	got := (&JSON{}).Render(nodes)

	assert.Equal(t, want, got)
	assert.True(t, json.Valid([]byte(got)), "output must be valid json")
}

func TestJSONRenderKeepsValueTypes(t *testing.T) {
	nodes := []Node{
		{Name: "int", Type: Changed, ValueBefore: 50.0, ValueAfter: 20.0},
		{Name: "float", Type: Changed, ValueBefore: 1.5, ValueAfter: 2.5},
		{Name: "string", Type: Changed, ValueBefore: "50", ValueAfter: "20"},
		{Name: "bool", Type: Changed, ValueBefore: false, ValueAfter: true},
		{Name: "null", Type: Changed, ValueBefore: nil, ValueAfter: nil},
		{Name: "emptyString", Type: Changed, ValueAfter: ""},
		{Name: "zero", Type: Changed, ValueAfter: 0.0},
	}

	var root map[string]any
	require.NoError(t, json.Unmarshal([]byte((&JSON{}).Render(nodes)), &root))

	children, ok := root["children"].([]any)
	require.True(t, ok)

	byKey := make(map[string]map[string]any, len(children))
	for _, raw := range children {
		child, ok := raw.(map[string]any)
		require.True(t, ok)

		key, ok := child["key"].(string)
		require.True(t, ok)

		byKey[key] = child
	}

	// Числа остаются числами, а не превращаются в строки.
	assert.Equal(t, float64(50), byKey["int"]["value1"])
	assert.Equal(t, float64(20), byKey["int"]["value2"])
	assert.Equal(t, 1.5, byKey["float"]["value1"])
	assert.Equal(t, "50", byKey["string"]["value1"])
	assert.Equal(t, "20", byKey["string"]["value2"])
	assert.Equal(t, false, byKey["bool"]["value1"])
	assert.Equal(t, true, byKey["bool"]["value2"])

	// omitempty у поля any пропускает только nil, поэтому false, 0 и ""
	// в вывод попадают, а отсутствующая сторона — нет.
	assert.NotContains(t, byKey["null"], "value1")
	assert.NotContains(t, byKey["null"], "value2")
	assert.NotContains(t, byKey["emptyString"], "value1")
	assert.Equal(t, "", byKey["emptyString"]["value2"])
	assert.NotContains(t, byKey["zero"], "value1")
	assert.Equal(t, float64(0), byKey["zero"]["value2"])
}

func TestJSONRenderNoNodes(t *testing.T) {
	assert.Equal(t, "{\n  \"key\": \"\",\n  \"type\": \"root\",\n  \"children\": []\n}", (&JSON{}).Render(nil))
}

func TestStatusJSON(t *testing.T) {
	assert.Equal(t, "added", statusJSON(Added))
	assert.Equal(t, "deleted", statusJSON(Deleted))
	assert.Equal(t, "changed", statusJSON(Changed))
	assert.Equal(t, "unchanged", statusJSON(Unchanged))
	assert.Equal(t, "nested", statusJSON(Nested), "nested обрабатывается до статуса")
}

func TestJSONRenderMarshalError(t *testing.T) {
	// yaml допускает .nan и .inf, которые json сериализовать не может —
	// на такой случай у вывода есть запасной вариант с ошибкой.
	got := (&JSON{}).Render([]Node{
		{Name: "broken", Type: Added, ValueAfter: math.NaN()},
	})

	assert.Contains(t, got, `"error"`)
	assert.Contains(t, got, "unsupported value")
}

func TestJSONRenderHasNoTrailingNewline(t *testing.T) {
	got := (&JSON{}).Render([]Node{
		{Name: "verbose", Type: Added, ValueAfter: true},
	})

	assert.False(t, strings.HasSuffix(got, "\n"))
}
