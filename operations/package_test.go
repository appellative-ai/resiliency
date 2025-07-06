package operations

import (
	"fmt"
	"github.com/behavioral-ai/resiliency/caseofficer"
	"os"
)

var (
	subDir            = "/operationstest/resource/"
	operatorsFileName = "logging-operators.json"
	appFileName       = "endpoint-config-old.json"
	endpointFileName  = "endpoint-config.json"
)

func readFile(fileName string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(dir + subDir + fileName)
}

func ExampleConfigureLogging() {
	err := ConfigureLogging(func() ([]byte, error) {
		return readFile(operatorsFileName)
	})
	fmt.Printf("test: ConfigureLogging(\"%v\") -> [err:%v]\n", subDir+operatorsFileName, err)

	//Output:
	//test: ConfigureLogging("/operationstest/resource/logging-operators.json") -> [err:<nil>]

}

func ExampleConfigureNetworks_Errors() {
	var appCfg map[string]map[string]string

	errs := ConfigureNetworks(appCfg, nil)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	errs = ConfigureNetworks(nil, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg = make(map[string]map[string]string)
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg["test"] = make(map[string]string)
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	//appCfg["test"] = "invalid file name"
	//errs = ConfigureNetworks(appCfg, readFile)
	//fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	//Output:
	//test: ConfigureNetworks() -> [network read function is nil]
	//test: ConfigureNetworks() -> [application config is nil or empty]
	//test: ConfigureNetworks() -> [application config is nil or empty]
	//test: ConfigureNetworks() -> [network file name is empty for case officer: test]

}

func ExampleConfigureNetworks() {
	//var appCfg map[string]string

	appCfg, err := ReadEndpointConfig(func() ([]byte, error) {
		return readFile(endpointFileName)
	})
	if err != nil {
		fmt.Printf("test: ReadEndpointConfig(\"%v\") -> [map:%v] [err:%v]\n", subDir+appFileName, len(appCfg), err)
	}

	errs := ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> [count:%v] [errs:%v]\n", len(errs), errs)

	a := opsAgent.Operative("core:common:agent/caseofficer/request/http/primary")
	if officer, ok := any(a).(caseofficer.Agent); ok {
		officer.Trace()
	}

	a = opsAgent.Operative("core:common:agent/caseofficer/request/http/secondary")
	if officer, ok := any(a).(caseofficer.Agent); ok {
		officer.Trace()
	}

	fmt.Printf("trace: Operations() -> %v\n", opsAgent.ex.List())

	//Output:
	//test: ConfigureNetworks() -> [count:0] [errs:[]]
	//trace: core:common:agent/caseofficer/request/http/primary -> test:resiliency:agent/rate-limiting/request/http
	//trace: core:common:agent/caseofficer/request/http/secondary -> test:resiliency:agent/routing/request/http
	//trace: Operations() -> [core:common:agent/caseofficer/request/http/primary core:common:agent/caseofficer/request/http/secondary]

}

func ExampleReadEndpointConfig() {
	cfg, err := ReadEndpointConfig(func() ([]byte, error) {
		return readFile(endpointFileName)
	})

	fmt.Printf("test: ReadEndpointConfig() -> %v [err:%v]\n", cfg, err)

	//Output:
	//test: ReadEndpointConfig() -> map[primary:map[network:network-config-primary.json pattern:/primary] secondary:map[network:network-config-secondary.json pattern:/secondary]] [err:<nil>]

}
