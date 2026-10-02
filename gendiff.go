package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code/formatters"
)

const DefaultFormat = "stylish"

// GenDiff сравнивает два файла и возвращает результат в выбранном формате.
// Пустой format означает DefaultFormat.
func GenDiff(path1, path2, format string) (string, error) {
	if format == "" {
		format = DefaultFormat
	}

	formatter, err := formatters.New(format)
	if err != nil {
		return "", err
	}

	before, err := loadFile(path1)
	if err != nil {
		return "", err
	}

	after, err := loadFile(path2)
	if err != nil {
		return "", err
	}

	return formatter.Render(formatters.BuildAST(before, after)), nil
}

func loadFile(path string) (map[string]any, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path %q: %w", path, err)
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", abs, err)
	}

	ext := strings.TrimPrefix(filepath.Ext(abs), ".")
	parsed, err := Parser(ext, data)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", abs, err)
	}

	m, ok := parsed.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data in %q is not an object", abs)
	}
	return m, nil
}
