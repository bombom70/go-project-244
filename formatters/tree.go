package formatters

import (
	"reflect"
	"slices"
)

type DiffType string

const (
	Added     DiffType = "added"
	Deleted   DiffType = "deleted"
	Unchanged DiffType = "unchanged"
	Changed   DiffType = "changed"
	Nested    DiffType = "nested"
)

type Node struct {
	Name        string
	ValueBefore any
	ValueAfter  any
	Type        DiffType
	Children    []Node
}

// BuildAST строит дерево различий двух объектов: ключи объединены и
// отсортированы, вложенные объекты разворачиваются в Children.
func BuildAST(data1, data2 map[string]any) []Node {
	keys := getKeys(data1, data2)
	nodes := make([]Node, 0, len(keys))

	for _, key := range keys {
		valueBefore, okBefore := data1[key]
		valueAfter, okAfter := data2[key]

		node := Node{Name: key, ValueBefore: valueBefore, ValueAfter: valueAfter}

		switch {
		case !okBefore:
			node.Type = Added
		case !okAfter:
			node.Type = Deleted
		case reflect.DeepEqual(valueBefore, valueAfter):
			node.Type = Unchanged
		default:
			bMap, bOK := valueBefore.(map[string]any)
			aMap, aOK := valueAfter.(map[string]any)
			if bOK && aOK {
				node.Type = Nested
				node.Children = BuildAST(bMap, aMap)
			} else {
				node.Type = Changed
			}
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func getKeys(data1, data2 map[string]any) []string {
	keys := make([]string, 0, len(data1)+len(data2))
	for k := range data1 {
		keys = append(keys, k)
	}
	for k := range data2 {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return slices.Compact(keys)
}
