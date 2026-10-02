package formatters

import (
	"fmt"
	"slices"
	"strings"
)

// Formatter — модуль вывода результата сравнения.
// Реализации лежат в отдельных файлах, выбор делает фабрика New.
type Formatter interface {
	Render(nodes []Node) string
}

// constructor создаёт форматтер для реестра.
type constructor func() Formatter

var registry = map[string]constructor{
	"json":    func() Formatter { return &JSON{} },
	"plain":   func() Formatter { return &Plain{} },
	"stylish": func() Formatter { return &Stylish{} },
}

// New — фабрика форматтеров: по имени отдаёт готовый Formatter.
func New(name string) (Formatter, error) {
	build, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown format: %q (available: %s)", name, strings.Join(Names(), ", "))
	}
	return build(), nil
}

// Names возвращает отсортированный список доступных форматов.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
