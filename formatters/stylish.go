package formatters

import (
	"fmt"
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
				indent(depth), n.Name, formatValue(n.ValueAfter)))

		case Deleted:
			lines = append(lines, fmt.Sprintf("%s- %s: %s",
				indent(depth), n.Name, formatValue(n.ValueBefore)))

		case Unchanged:
			lines = append(lines, fmt.Sprintf("%s  %s: %s",
				indent(depth), n.Name, formatValue(n.ValueBefore)))

		case Changed:
			lines = append(lines,
				fmt.Sprintf("%s- %s: %s", indent(depth), n.Name, formatValue(n.ValueBefore)),
				fmt.Sprintf("%s+ %s: %s", indent(depth), n.Name, formatValue(n.ValueAfter)),
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

// formatValue сериализует значение узла. Составные значения (map, слайс)
// не разворачиваются в дерево, а обозначаются как [complex value] —
// на их месте в before/after может стоять объект целиком.
func formatValue(v any) string {
	switch val := v.(type) {
	case map[string]any, []any:
		return "[complex value]"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
