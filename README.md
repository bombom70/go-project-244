# Вычислитель отличий на Go

[![hexlet-check](https://github.com/bombom70/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/bombom70/go-project-244/actions)

Консольная утилита для сравнения вложенных структур (JSON, YAML)

Учебный проект Хекслета: https://ru.hexlet.io/programs/go


## Стек

- Go

## Установка

<!-- Опишите установку: клонирование, зависимости, переменные окружения -->

```bash
git clone https://github.com/bombom70/go-project-244.git
cd go-project-244
make build
```

## Использование

<!-- Добавьте примеры запуска и запись asciinema — именно это смотрит работодатель -->

Утилита принимает пути к двум файлам и флаг `--format` (`-f`) с форматом вывода.
Формат по умолчанию — `stylish`.

```bash
./bin/gendiff testdata/fixture/beforeTree.json testdata/fixture/afterTree.json

./bin/gendiff --format stylish testdata/fixture/beforeTree.json testdata/fixture/afterTree.json

./bin/gendiff --format plain testdata/fixture/beforeTree.json testdata/fixture/afterTree.json

./bin/gendiff --format json testdata/fixture/beforeTree.json testdata/fixture/afterTree.json

./bin/gendiff --format plain testdata/fixture/before.yaml testdata/fixture/after.yaml
```

### Форматы

| Формат   | Описание                                                                            |
|----------|-------------------------------------------------------------------------------------|
| `stylish` | Вложенный вид с отступами: `+` — добавлено, `-` — удалено, `-`/`+` — изменено         |
| `plain`   | Плоский список, по одной строке на каждое изменённое свойство                        |
| `json`    | Структурированный вывод: статус и оба значения каждого свойства                      |

В формате `plain` имя свойства выводится с полным путём от корня
(`common.setting6.ops`), строки оборачиваются в одинарные кавычки
(`'blah blah'`), числа, `true`, `false` и `null` — как есть, составные
значения — как `[complex value]`:

```
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
```

В формате `json` каждое свойство описывается статусом
(`unchanged`, `added`, `removed`, `changed`) и обоими значениями, а вложенные
объекты остаются вложенными объектами. Значения сохраняют свои типы, поэтому
вывод можно передать другой программе:

```json
{
  "common": {
    "follow": {
      "status": "added",
      "beforeValue": null,
      "afterValue": false
    },
    "setting2": {
      "status": "removed",
      "beforeValue": 200,
      "afterValue": null
    },
    "setting6": {
      "doge": {
        "wow": {
          "status": "changed",
          "beforeValue": "",
          "afterValue": "so much"
        }
      },
      "ops": {
        "status": "added",
        "beforeValue": null,
        "afterValue": "vops"
      }
    }
  },
  "group2": {
    "status": "removed",
    "beforeValue": {
      "abc": 12345,
      "deep": {
        "id": 45
      }
    },
    "afterValue": null
  }
}
```

### Пример работы пакета

`GenDiff` сам читает и разбирает файлы (JSON, YAML, INI) и печатает
результат в выбранном формате:

```go
package main

import (
	"fmt"
	"log"

	"code"
)

func main() {
	out, err := code.GenDiff(
		"testdata/fixture/beforeTree.json",
		"testdata/fixture/afterTree.json",
		"plain",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
```

Формат можно выбрать и через параметр функции — пустая строка означает
формат по умолчанию (`stylish`).

Если нужно отдельно построить дерево различий и отрисовать его своим
форматтером, используются модуль `formatters` и интерфейс `Formatter`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"code/formatters"
)

func main() {
	formatter, err := formatters.New("plain")
	if err != nil {
		log.Fatal(err)
	}

	before := load("testdata/fixture/beforeTree.json")
	after := load("testdata/fixture/afterTree.json")

	fmt.Println(formatter.Render(formatters.BuildAST(before, after)))
}

func load(path string) map[string]any {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		log.Fatal(err)
	}

	return out
}
```

Свой форматтер достаточно реализовать как метод `Render([]Node) string` и
добавить в реестр `formatters/formatters.go`:

```go
type Formatter interface {
	Render(nodes []Node) string
}
```

## Устройство проекта

```
cmd/gendiff/        точка входа CLI
formatters/
  formatters.go     интерфейс Formatter и фабрика New — выбор формата
  tree.go           Node и BuildAST — дерево различий
  stylish.go        форматтер stylish
  plain.go          форматтер plain
  json.go           форматтер json
gendiff.go          GenDiff: чтение файлов и выбор формата
parser.go           разбор JSON, YAML, INI
testdata/fixture/   файлы для сравнения и эталонные результаты
```

Форматы не зашиты в `GenDiff`: он получает готовый `Formatter` из фабрики
`formatters.New`, поэтому новый формат добавляется одним файлом в
`formatters` без изменений в коде сравнения.

## Разработка

```bash
make test           тесты
make test-coverage  тесты с отчётом о покрытии (минимум 70%)
make lint           статический анализ golangci-lint
make build          сборка бинарника в bin/gendiff
```

---

<details>
<summary>Автоматические тесты Хекслета</summary>

Тесты запускаются на каждый коммит. За запуск отвечает файл `.github/workflows/hexlet-check.yml` — не удаляйте и не переименовывайте ни его, ни репозиторий.

</details>

## О Хекслете

[Хекслет](https://ru.hexlet.io/) — школа программирования: авторские программы обучения с практикой, поддержкой наставников и реальными проектами, которые остаются в резюме. Этот репозиторий — один из таких проектов.
