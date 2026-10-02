package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

const (
	indentSize = 4 // один уровень вложенности
	shiftLeft  = 2 // смещение влево для + / -
)

type Stylish struct{}

func (s *Stylish) Render(nodes []Node) string {
	return "{\n" + renderNodes(nodes, 1) + "\n}"
}

func renderNodes(nodes []Node, depth int) string {
	lines := make([]string, 0, len(nodes))

	for _, n := range nodes {
		switch n.Type {
		case Added:
			lines = append(lines, fmt.Sprintf("%s+ %s: %s",
				indent(depth), n.Name, formatValue(n.ValueAfter, depth)))

		case Deleted:
			lines = append(lines, fmt.Sprintf("%s- %s: %s",
				indent(depth), n.Name, formatValue(n.ValueBefore, depth)))

		case Unchanged:
			lines = append(lines, fmt.Sprintf("%s  %s: %s",
				indent(depth), n.Name, formatValue(n.ValueBefore, depth)))

		case Changed:
			lines = append(lines,
				fmt.Sprintf("%s- %s: %s", indent(depth), n.Name, formatValue(n.ValueBefore, depth)),
				fmt.Sprintf("%s+ %s: %s", indent(depth), n.Name, formatValue(n.ValueAfter, depth)),
			)

		case Nested:
			lines = append(lines, fmt.Sprintf("%s  %s: {\n%s\n%s}",
				indent(depth), n.Name,
				renderNodes(n.Children, depth+1),
				closeIndent(depth),
			))
		}
	}

	return strings.Join(lines, "\n")
}

// indent возвращает пробелы для узла дифа на глубине depth.
// Формула: depth*indentSize - shiftLeft
func indent(depth int) string {
	n := depth*indentSize - shiftLeft
	if n < 0 {
		n = 0
	}
	return strings.Repeat(" ", n)
}

// closeIndent возвращает пробелы для закрывающей скобки вложенного объекта:
// она стоит в той же колонке, что и имя объекта, то есть на indent(depth)+2.
func closeIndent(depth int) string {
	return indent(depth) + strings.Repeat(" ", 2)
}

// formatValue сериализует значение узла на глубине depth.
// Объект, который не удалось разложить на поддерево (в before/after он стоит
// целиком), печатается как вложенный блок с отсортированными ключами: скобки
// на колонке depth*indentSize, свойства на depth*indentSize+indentSize.
// Массив и прочие составные значения обозначаются как [complex value].
func formatValue(v any, depth int) string {
	obj, ok := v.(map[string]any)
	if !ok {
		return formatScalar(v)
	}

	lines := make([]string, 0, len(obj))
	for _, key := range slices.Sorted(maps.Keys(obj)) {
		lines = append(lines, fmt.Sprintf("%s%s: %s",
			strings.Repeat(" ", depth*indentSize+indentSize),
			key,
			formatValue(obj[key], depth+1),
		))
	}

	if len(lines) == 0 {
		return "{}"
	}

	return fmt.Sprintf("{\n%s\n%s}", strings.Join(lines, "\n"),
		strings.Repeat(" ", depth*indentSize))
}

// formatScalar печатает значение, которое не является объектом: null как
// null, массивы и прочие составные значения — как [complex value],
// остальное — через %v.
func formatScalar(v any) string {
	switch v.(type) {
	case []any:
		return "[complex value]"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}
