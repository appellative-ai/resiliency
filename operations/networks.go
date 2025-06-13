package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/resiliency/caseofficer"
	//"github.com/behavioral-ai/resiliency/endpoint"
)

func configureNetworks(appCfg map[string]string, read func(fileName string) ([]byte, error)) (errs []error) {
	if read == nil {
		return []error{errors.New("network read function is nil")}
	}
	if len(appCfg) == 0 {
		return []error{errors.New("application config is nil or empty")}
	}

	//var result = make([]error, len(appCfg)*2)
	//var wg sync.WaitGroup
	//var i int

	for k, v := range appCfg {
		if v == "" {
			errs = append(errs, errors.New(fmt.Sprintf("file name is empty for case officer: %v", k)))
			continue
		}
		officer, err1 := validateOfficerType(k)
		if err1 != nil {
			errs = append(errs, err1)
			continue
		}
		netCfg, err := buildNetworkConfig(v, read)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		chain, errs1 := officer.BuildNetwork(netCfg)
		if errs1 != nil {
			errs = append(errs, errs1...)
			continue
		}
		err = buildEndpoint(officer.Name(), chain)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		//wg.Add(1)
		//buildEndpoint(officer, v, read, &result[i])
	}
	//wg.Wait()
	// Need to create
	return packErrors(errs)
}

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
