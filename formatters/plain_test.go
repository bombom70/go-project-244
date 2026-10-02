package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainRender(t *testing.T) {
	tests := []struct {
		name  string
		nodes []Node
		want  string
	}{
		{
			name:  "no nodes produce empty output",
			nodes: nil,
			want:  "",
		},
		{
			name: "added removed and updated properties",
			nodes: []Node{
				{Name: "follow", Type: Added, ValueAfter: false},
				{Name: "timeout", Type: Changed, ValueBefore: 50.0, ValueAfter: 20.0},
				{Name: "proxy", Type: Deleted, ValueBefore: "123.234.53.22"},
			},
			want: "Property 'follow' was added with value: false\n" +
				"Property 'timeout' was updated. From 50 to 20\n" +
				"Property 'proxy' was removed",
		},
		{
			name: "unchanged properties are skipped",
			nodes: []Node{
				{Name: "host", Type: Unchanged, ValueBefore: "hexlet.io", ValueAfter: "hexlet.io"},
				{Name: "verbose", Type: Added, ValueAfter: true},
			},
			want: "Property 'verbose' was added with value: true",
		},
		{
			name: "strings are quoted and null is printed as is",
			nodes: []Node{
				{Name: "setting3", Type: Changed, ValueBefore: true, ValueAfter: nil},
				{Name: "setting4", Type: Added, ValueAfter: "blah blah"},
			},
			want: "Property 'setting3' was updated. From true to null\n" +
				"Property 'setting4' was added with value: 'blah blah'",
		},
		{
			name: "complex values are collapsed",
			nodes: []Node{
				{Name: "setting5", Type: Added, ValueAfter: map[string]any{"key5": "value5"}},
				{Name: "hosts", Type: Added, ValueAfter: []any{"a", "b"}},
				{Name: "nest", Type: Changed, ValueBefore: map[string]any{"key": "value"}, ValueAfter: "str"},
				{Name: "group2", Type: Deleted, ValueBefore: map[string]any{"abc": 12345.0}},
			},
			want: "Property 'setting5' was added with value: [complex value]\n" +
				"Property 'hosts' was added with value: [complex value]\n" +
				"Property 'nest' was updated. From [complex value] to 'str'\n" +
				"Property 'group2' was removed",
		},
		{
			name: "nested properties are reported with the full path",
			nodes: []Node{
				{Name: "common", Type: Nested, Children: []Node{
					{Name: "setting1", Type: Unchanged, ValueBefore: "Value 1", ValueAfter: "Value 1"},
					{Name: "setting6", Type: Nested, Children: []Node{
						{Name: "doge", Type: Nested, Children: []Node{
							{Name: "wow", Type: Changed, ValueBefore: "", ValueAfter: "so much"},
						}},
						{Name: "ops", Type: Added, ValueAfter: "vops"},
					}},
				}},
			},
			want: "Property 'common.setting6.doge.wow' was updated. From '' to 'so much'\n" +
				"Property 'common.setting6.ops' was added with value: 'vops'",
		},
		{
			name: "nested branch without changes adds no lines",
			nodes: []Node{
				{Name: "untouched", Type: Nested, Children: []Node{
					{Name: "same", Type: Unchanged, ValueBefore: 1.0, ValueAfter: 1.0},
				}},
				{Name: "host", Type: Added, ValueAfter: "hexlet.io"},
			},
			want: "Property 'host' was added with value: 'hexlet.io'",
		},
	}

	plain := &Plain{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, plain.Render(tt.nodes))
		})
	}
}

func TestPlainRenderHasNoTrailingNewline(t *testing.T) {
	got := (&Plain{}).Render([]Node{
		{Name: "verbose", Type: Added, ValueAfter: true},
		{Name: "proxy", Type: Deleted, ValueBefore: "123.234.53.22"},
	})

	assert.Equal(t, "Property 'verbose' was added with value: true\nProperty 'proxy' was removed", got)
}
