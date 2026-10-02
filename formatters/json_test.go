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
  "added": {
    "status": "added",
    "beforeValue": null,
    "afterValue": true
  },
  "changed": {
    "status": "changed",
    "beforeValue": "bas",
    "afterValue": "bars"
  },
  "removed": {
    "status": "removed",
    "beforeValue": 200,
    "afterValue": null
  },
  "same": {
    "status": "unchanged",
    "beforeValue": "hexlet.io",
    "afterValue": "hexlet.io"
  }
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
  "common": {
    "setting6": {
      "doge": {
        "wow": {
          "status": "changed",
          "beforeValue": "",
          "afterValue": "so much"
        }
      },
      "key": {
        "status": "unchanged",
        "beforeValue": "value",
        "afterValue": "value"
      }
    }
  },
  "empty": {},
  "group2": {
    "status": "removed",
    "beforeValue": {
      "abc": 12345
    },
    "afterValue": null
  },
  "group3": {
    "status": "added",
    "beforeValue": null,
    "afterValue": {
      "fee": 100500
    }
  },
  "hosts": {
    "status": "added",
    "beforeValue": null,
    "afterValue": [
      "a",
      "b"
    ]
  }
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
	}

	var got map[string]map[string]any
	require.NoError(t, json.Unmarshal([]byte((&JSON{}).Render(nodes)), &got))

	// Числа остаются числами, а не превращаются в строки.
	assert.Equal(t, float64(50), got["int"]["beforeValue"])
	assert.Equal(t, float64(20), got["int"]["afterValue"])
	assert.Equal(t, 1.5, got["float"]["beforeValue"])
	assert.Equal(t, "50", got["string"]["beforeValue"])
	assert.Equal(t, "20", got["string"]["afterValue"])
	assert.Equal(t, false, got["bool"]["beforeValue"])
	assert.Equal(t, true, got["bool"]["afterValue"])
	assert.Nil(t, got["null"]["beforeValue"])
}

func TestJSONRenderNoNodes(t *testing.T) {
	assert.Equal(t, "{}", (&JSON{}).Render(nil))
}

func TestStatusJSON(t *testing.T) {
	// Внутри дерева удалённое свойство помечено как deleted, а в вывод
	// уходит как removed — так же, как в формате plain.
	assert.Equal(t, "added", statusJSON(Added))
	assert.Equal(t, "removed", statusJSON(Deleted))
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
