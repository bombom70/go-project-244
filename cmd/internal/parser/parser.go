package parser

import (
	"os"
	"path"
)

func Parse() {

}

func readFile(currentPath string) ([]byte, error) {
	rootPath := path.Join(currentPath)
	data, err := os.ReadFile(rootPath)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
