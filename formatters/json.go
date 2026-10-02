package formatters

import (
	"encoding/json"
	"fmt"
)

// jsonIndent — отступ в два пробела, как у json.MarshalIndent по умолчанию.
const jsonIndent = "  "

type JSON struct{}

// jsonDiff — представление одного изменённого свойства в формате json.
type jsonDiff struct {
	Status      string `json:"status"`
	BeforeValue any    `json:"beforeValue"`
	AfterValue  any    `json:"afterValue"`
}

func (j *JSON) Render(nodes []Node) string {
	// Значения приходят из разбора json/yaml, поэтому сериализовать их
	// нечем: MarshalIndent падает только на каналах и функциях. На
	// нештатный случай отдаём объект с ошибкой, а не пустую строку.
	out, err := json.MarshalIndent(toJSONValue(nodes), "", jsonIndent)
	if err != nil {
		return fmt.Sprintf("{%q: %q}", "error", err.Error())
	}

	return string(out)
}

// toJSONValue превращает дерево различий в объекты, готовые к сериализации:
// вложенный объект остаётся под-объектом без статуса, всё остальное
// превращается в пару before/after со статусом.
func toJSONValue(nodes []Node) map[string]any {
	out := make(map[string]any, len(nodes))

	for _, n := range nodes {
		if n.Type == Nested {
			out[n.Name] = toJSONValue(n.Children)
			continue
		}

		out[n.Name] = jsonDiff{
			Status:      statusJSON(n.Type),
			BeforeValue: n.ValueBefore,
			AfterValue:  n.ValueAfter,
		}
	}

	return out
}

// statusJSON переводит DiffType в статус формата json. Удалённое свойство
// в дереве помечено как deleted, а в выводе — как removed, как и в plain.
func statusJSON(t DiffType) string {
	switch t {
	case Added:
		return "added"
	case Deleted:
		return "removed"
	case Changed:
		return "changed"
	case Unchanged:
		return "unchanged"
	default:
		return string(t)
	}
}
