package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/resiliency/caseofficer"
)

func buildNetworkConfig(fileName string, read func(fileName string) ([]byte, error)) (map[string]map[string]string, error) {
	var buf []byte
	var err error
	var appCfg []map[string]string

	buf, err = read(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &appCfg)
	if err != nil {
		return nil, err
	}
	return shapeNetworkConfig(appCfg), nil
}

func shapeNetworkConfig(cfg []map[string]string) map[string]map[string]string {
	newCfg := make(map[string]map[string]string)
	for _, m := range cfg {
		newCfg[m[RoleKey]] = m
		delete(m, RoleKey)
	}
	return newCfg
}

func validateOfficerType(name string) (caseofficer.Agent, error) {
	if name == "" {
		return nil, errors.New(fmt.Sprintf("case officer name is empty"))
	}
	agent := repository.Agent(name)
	if agent == nil {
		return nil, errors.New(fmt.Sprintf("agent lookup is nil for case officer: %v", name))
	}
	officer, ok := any(agent).(caseofficer.Agent)
	if !ok {
		return nil, errors.New(fmt.Sprintf("agent is not of type caseofficer.Agent for case officer: %v", name))
	}
	return officer, nil
}

func packErrors(errs []error) []error {
	var result []error
	for _, err := range errs {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}

//([]byte, error), err *error)
/*{
	defer wg.Done()
	var networkCfg []map[string]string

	buf, err2 := read(fileName)
	if err2 != nil {
		*err = err2
		return
	}
	err2 = json.Unmarshal(buf, networkCfg)
	if err2 != nil {
		*err = err2
		return
	}
}(officer, fileName, read, &result[i])

*/

/*
func validateOfficerConfig(caseOfficer string, m map[string]string) (name string, fileName string, err error) {
	var ok bool
	name, ok = m[NameKey]
	if name == "" || !ok {
		return "", "", errors.New(fmt.Sprintf("name is empty for case officer: %v", caseOfficer))
	}
	fileName, ok = m[PathKey]
	if fileName == "" || !ok {
		return "", "", errors.New(fmt.Sprintf("file name is empty for case officer: %v", caseOfficer))
	}
	return
}


*/

/*
func parseOfficerConfig(s string) map[string]string {
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

*/
