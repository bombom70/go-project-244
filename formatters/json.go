package formatters

import (
	"encoding/json"
	"fmt"
)

// jsonIndent — отступ в два пробела, как у json.MarshalIndent по умолчанию.
const jsonIndent = "  "

type JSON struct{}

// jsonNode — представление одного узла дерева различий в формате json.
//
// Значения хранятся в value1 (было) и value2 (стало) и печатаются только
// когда они есть: omitempty у поля типа any пропускает лишь nil, поэтому
// false, 0 и "" в вывод попадают.
//
// Children — указатель на срез, а не сам срез: omitempty у nil-указателя
// убирает поле совсем, а непустой срез под ним печатается даже пустым.
// Так у корня и у вложенного объекта children есть всегда, в том числе
// когда поддерево пустое, и форма вывода не зависит от наличия данных.
type jsonNode struct {
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Value1   any         `json:"value1,omitempty"`
	Value2   any         `json:"value2,omitempty"`
	Children *[]jsonNode `json:"children,omitempty"`
}

func (j *JSON) Render(nodes []Node) string {
	root := jsonNode{Key: "", Type: "root", Children: children(nodes)}

	// Значения приходят из разбора json/yaml, поэтому сериализовать их
	// нечем: MarshalIndent падает только на каналах и функциях. На
	// нештатный случай отдаём объект с ошибкой, а не пустую строку.
	out, err := json.MarshalIndent(root, "", jsonIndent)
	if err != nil {
		return fmt.Sprintf("{%q: %q}", "error", err.Error())
	}

	return string(out)
}

// children превращает дерево различий в готовые к сериализации узлы.
// Срез никогда не nil, поэтому вложенный объект без детей даёт [].
func children(nodes []Node) *[]jsonNode {
	out := make([]jsonNode, 0, len(nodes))

	for _, n := range nodes {
		node := jsonNode{Key: n.Name, Type: statusJSON(n.Type)}

		switch n.Type {
		case Added:
			node.Value2 = n.ValueAfter
		case Deleted:
			node.Value1 = n.ValueBefore
		case Changed:
			node.Value1 = n.ValueBefore
			node.Value2 = n.ValueAfter
		case Unchanged:
			node.Value1 = n.ValueBefore
		case Nested:
			node.Children = children(n.Children)
		}

		out = append(out, node)
	}

	return &out
}

// statusJSON переводит DiffType в статус формата json. Названия совпадают
// с DiffType, в том числе deleted у удалённого свойства.
func statusJSON(t DiffType) string {
	return string(t)
}
