package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name     string
		path1    string
		path2    string
		format   string
		wantFile string
	}{
		{
			name:     "stylish with flat json",
			path1:    "./testdata/fixture/before.json",
			path2:    "./testdata/fixture/after.json",
			format:   "stylish",
			wantFile: "./testdata/fixture/stylishFlat.txt",
		},
		{
			name:     "stylish with nested json",
			path1:    "./testdata/fixture/beforeTree.json",
			path2:    "./testdata/fixture/afterTree.json",
			format:   "stylish",
			wantFile: "./testdata/fixture/stylishTree.txt",
		},
		{
			name:     "plain with nested json",
			path1:    "./testdata/fixture/beforeTree.json",
			path2:    "./testdata/fixture/afterTree.json",
			format:   "plain",
			wantFile: "./testdata/fixture/plainTree.txt",
		},
		{
			name:     "json with nested json",
			path1:    "./testdata/fixture/beforeTree.json",
			path2:    "./testdata/fixture/afterTree.json",
			format:   "json",
			wantFile: "./testdata/fixture/jsonTree.json",
		},
		{
			name:     "json with yaml",
			path1:    "./testdata/fixture/before.yaml",
			path2:    "./testdata/fixture/after.yaml",
			format:   "json",
			wantFile: "./testdata/fixture/jsonFlat.json",
		},
		{
			name:     "stylish with yaml gives the same result as json",
			path1:    "./testdata/fixture/before.yaml",
			path2:    "./testdata/fixture/after.yaml",
			format:   "stylish",
			wantFile: "./testdata/fixture/stylishFlat.txt",
		},
		{
			name:     "empty format falls back to stylish",
			path1:    "./testdata/fixture/before.json",
			path2:    "./testdata/fixture/after.json",
			format:   "",
			wantFile: "./testdata/fixture/stylishFlat.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(tt.path1, tt.path2, tt.format)
			require.NoError(t, err)

			wantBytes, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)

			assert.Equal(t, strings.TrimSpace(string(wantBytes)), strings.TrimSpace(got))
		})
	}
}

// TestGenDiffHexletFixtures сверяет вывод с эталонными файлами из задания
// Хекслета. Оба файла сравнения есть в json и в yaml, и результат должен
// совпадать с одним и тем же эталоном в любом из трёх форматов.
func TestGenDiffHexletFixtures(t *testing.T) {
	tests := []struct {
		format   string
		wantFile string
	}{
		{format: "stylish", wantFile: "./testdata/fixture/result_stylish.txt"},
		{format: "plain", wantFile: "./testdata/fixture/result_plain.txt"},
		{format: "json", wantFile: "./testdata/fixture/result_json.json"},
	}
	extensions := []string{"json", "yml"}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			wantBytes, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)
			want := strings.TrimSpace(string(wantBytes))

			for _, ext := range extensions {
				got, err := GenDiff(
					"./testdata/fixture/file1."+ext,
					"./testdata/fixture/file2."+ext,
					tt.format,
				)
				require.NoError(t, err)

				assert.Equal(t, want, strings.TrimSpace(got), "input extension .%s", ext)
			}
		})
	}
}

func TestGenDiffErrors(t *testing.T) {
	dir := t.TempDir()

	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
		return path
	}

	notObject := write("array.json", `[1, 2, 3]`)
	unknownExt := write("data.txt", "timeout=20")
	broken := write("broken.json", `{"timeout": `)
	ok := write("ok.json", `{"timeout": 20}`)

	tests := []struct {
		name        string
		path1       string
		path2       string
		format      string
		errContains string
	}{
		{
			name:        "unknown format",
			path1:       ok,
			path2:       ok,
			format:      "yaml",
			errContains: `unknown format: "yaml"`,
		},
		{
			name:        "format is case sensitive",
			path1:       ok,
			path2:       ok,
			format:      "JSON",
			errContains: `unknown format: "JSON"`,
		},
		{
			name:        "missing first file",
			path1:       filepath.Join(dir, "nope.json"),
			path2:       ok,
			format:      "plain",
			errContains: "no such file or directory",
		},
		{
			name:        "missing second file",
			path1:       ok,
			path2:       filepath.Join(dir, "nope.json"),
			format:      "plain",
			errContains: "no such file or directory",
		},
		{
			name:        "unknown extension",
			path1:       unknownExt,
			path2:       ok,
			format:      "plain",
			errContains: `unknown parser type: "txt"`,
		},
		{
			name:        "broken content",
			path1:       broken,
			path2:       ok,
			format:      "plain",
			errContains: "parse",
		},
		{
			name:        "root is not an object",
			path1:       notObject,
			path2:       ok,
			format:      "plain",
			errContains: "is not an object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(tt.path1, tt.path2, tt.format)
			require.Error(t, err)
			assert.Empty(t, got)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}
