package util

import (
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

// returns key, value for all environment variables starting with prefix
func environment(prefix string) map[string]string {
	env := make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		if strings.HasPrefix(pair[0], prefix) {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

func customFuncMap() template.FuncMap {
	var functionMap = map[string]interface{}{"environment": environment}
	return template.FuncMap(functionMap)
}

func Parse(s string) (*template.Template, error) {
	return template.New("").Funcs(sprig.TxtFuncMap()).Funcs(customFuncMap()).Parse(s)
}

func ReadEnv() (env map[string]string) {
	env = make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		env[pair[0]] = pair[1]
	}
	return
}
