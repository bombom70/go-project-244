package code

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func GenDiff(filepath1, filepath2, format string) string {
	fileData1, err := Parser(filepath1)
	if err != nil {
		return fmt.Sprintf("Не удалось распарсить файл %s: %v", filepath1, err)
	}

	fileData2, err := Parser(filepath2)
	if err != nil {
		return fmt.Sprintf("Не удалось распарсить файл %s: %d", filepath2, err)
	}

	nodes := buildAst(fileData1, fileData2)
	res := render(nodes)
	fmt.Println(res)
	return "Path"
}

type DiffType string

const (
	Added     DiffType = "added"
	Deleted   DiffType = "deleted"
	Unchanged DiffType = "unchanged"
	Changed   DiffType = "changeInside"
	Space     int      = 2
)

type Node struct {
	Name        string
	ValueBefore any
	ValueAfter  any
	Type        DiffType
	Children    []Node
}

func buildAst(data1, data2 map[string]any) []Node {
	ast := []Node{}
	allKeys := getKeys(data1, data2)
	keys := make([]string, 0, len(allKeys))
	for _, k := range allKeys {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		var node Node
		valueBefore, okBefore := data1[key]
		valueAfter, okAfter := data2[key]

		if !okBefore {
			node.Type = Added
		} else if !okAfter {
			node.Type = Deleted
		} else if valueAfter == valueBefore {
			node.Type = Unchanged
		} else {
			node.Type = Changed
		}

		node.Name = key
		node.ValueBefore = valueBefore
		node.ValueAfter = valueAfter
		ast = append(ast, node)
	}
	// fmt.Println(ast)

	return ast
}

func render(nodes []Node) string {
	var strs []string

	for _, n := range nodes {
		var s string
		switch n.Type {
		case Added:
			s = fmt.Sprintf("%s+ %s: %v", strings.Repeat(" ", Space), n.Name, n.ValueAfter)
		case Deleted:
			s = fmt.Sprintf("%s- %s: %v", strings.Repeat(" ", Space), n.Name, n.ValueBefore)
		case Unchanged:
			s = fmt.Sprintf("%s%s: %v", strings.Repeat(" ", Space*2), n.Name, n.ValueBefore)
		case Changed:
			s = fmt.Sprintf("%s- %s: %v\n%s+ %s: %v", strings.Repeat(" ", Space), n.Name, n.ValueBefore, strings.Repeat(" ", 2), n.Name, n.ValueAfter)
		default:
			// Кидать ошибку
		}
		strs = append(strs, s)
	}

	return fmt.Sprintf("{\n%v\n}", strings.Join(strs, "\n"))
}

func getKeys(data1, data2 map[string]any) []string {
	var keys []string

	for key := range data1 {
		if !slices.Contains(keys, key) {
			keys = append(keys, key)
		}
	}

	for key := range data2 {
		if !slices.Contains(keys, key) {
			keys = append(keys, key)
		}
	}

	return keys
}
