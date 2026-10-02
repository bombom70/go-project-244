package code

import (
	"fmt"

	"encoding/json"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

func Parser(t string, data []byte) (any, error) {
	f, ok := mapping[t]
	if !ok {
		return nil, fmt.Errorf("unknown parser type: %q", t)
	}
	return f(data)
}

var mapping = map[string]func([]byte) (any, error){
	"json": func(data []byte) (any, error) {
		var out any
		err := json.Unmarshal(data, &out)
		return out, err
	},
	"yaml": func(data []byte) (any, error) {
		var out any
		err := yaml.Unmarshal(data, &out)
		return out, err
	},
	"yml": func(data []byte) (any, error) {
		var out any
		err := yaml.Unmarshal(data, &out)
		return out, err
	},
	"ini": func(data []byte) (any, error) {
		return ini.Load(data)
	},
}
