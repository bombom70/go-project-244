package formatters

import (
	"fmt"
	"strings"
)

type Plain struct{}

func (p *Plain) Render(nodes []Node) string {
	return strings.Join(renderPlain(nodes, ""), "\n")
}

// renderPlain обходит дерево и собирает строки отчёта.
// prefix — путь до текущего узла от корня, чтобы в имени показывать
// весь путь (например common.setting6.ops), а не только имя родителя.
func renderPlain(nodes []Node, prefix string) []string {
	lines := make([]string, 0, len(nodes))

	for _, n := range nodes {
		name := prefix + n.Name

		switch n.Type {
		case Added:
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %s",
				name, plainValue(n.ValueAfter)))

		case Deleted:
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", name))

		case Changed:
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
				name, plainValue(n.ValueBefore), plainValue(n.ValueAfter)))

		case Nested:
			lines = append(lines, renderPlain(n.Children, name+".")...)

		case Unchanged:
			continue
		}
	}

	return lines
}

// plainValue сериализует значение для плоского вывода: строки в одинарных
// кавычках, составные значения (map, слайс) — как [complex value],
// числа, булевы значения и null — как есть.
func plainValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", val)
	case map[string]any, []any:
		return "[complex value]"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
