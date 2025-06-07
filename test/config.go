package test

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	nameKey = "name"
	pathKey = "@path"
)

type AppConfig struct {
	app map[string]map[string]string
}

func NewAppConfig(m map[string]string) *AppConfig {
	a := new(AppConfig)
	a.app = make(map[string]map[string]string)
	if m != nil {
		for k, v := range m {
			a.app[k] = ParseValue(v)
		}
	}
	return a
}

func (a *AppConfig) Name(k string) (string, bool) {
	if v, ok := a.app[k]; ok {
		return v[nameKey], ok
	}
	return "", false
}

func (a *AppConfig) Path(k string) (string, bool) {
	if v, ok := a.app[k]; ok {
		return v[pathKey], ok
	}
	return "", false
}

func ReadConfig[T any](path string) (t T, err error) {
	var buf []byte
	var dir string

	dir, err = os.Getwd()
	if err != nil {
		return
	}
	buf, err = os.ReadFile(dir + path)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(buf, &t)
	if err != nil {
		return t, err
	}
	return
}

func ParseValue(s string) map[string]string {
	var m = make(map[string]string)

	tokens := strings.Split(s, ",")
	for _, t := range tokens {
		pairs := strings.Split(t, "=")
		if len(pairs) < 2 || pairs[1] == "" {
			continue
		}
		m[pairs[0]] = pairs[1]
	}
	return m
}
