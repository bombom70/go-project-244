package code

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// TODO: Думаю что надо так парсить для расширени других форматов
// type ParserS struct{}
// func (p *ParserS) json(data []byte) (map[string]any, error) {
// 	var result map[string]any
// 	err := json.Unmarshal(data, &result)
// 	return result, err
// }

func Parser(filePath string) (map[string]any, error) {
	data, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	return result, err
}

func ReadFile(filePath string) ([]byte, error) {
	rootPath, err := makeFilePath(filePath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(rootPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func makeFilePath(filePath string) (string, error) {
	if filepath.IsAbs(filePath) {
		return filepath.Clean(filePath), nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Clean(filepath.Join(dir, filePath)), nil
}
