package network

import (
	"encoding/json"
	"errors"
)

const (
	roleKey = "role"
)

func BuildConfig(fileName string, read func(fileName string) ([]byte, error)) (map[string]map[string]string, error) {
	var buf []byte
	var err error
	var appCfg []map[string]string

	if read == nil {
		return nil, errors.New("network read function is nil")
	}
	if fileName == "" {
		return nil, errors.New("application config is nil or empty")
	}
	buf, err = read(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &appCfg)
	if err != nil {
		return nil, err
	}
	return ShapeConfig(roleKey, appCfg), nil
}

func ShapeConfig(mapKey string, cfg []map[string]string) map[string]map[string]string {
	newCfg := make(map[string]map[string]string)
	for _, m := range cfg {
		newCfg[m[mapKey]] = m
		delete(m, mapKey)
	}
	return newCfg
}

/*
func packErrors(errs []error) []error {
	var result []error
	for _, err := range errs {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}

*/
