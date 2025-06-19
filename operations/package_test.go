package operations

import (
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/resiliency/caseofficer"
	"github.com/behavioral-ai/resiliency/module"
	"os"
)

var (
	subDir            = "/operationstest/resource/"
	operatorsFileName = "logging-operators.json"
	appFileName       = "app-config.json"
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
	var appCfg map[string]string

	errs := ConfigureNetworks(appCfg, nil)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	errs = ConfigureNetworks(nil, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg = make(map[string]string)
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg["test"] = ""
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg["test"] = "invalid file name"
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	//Output:
	//test: ConfigureNetworks() -> [network read function is nil]
	//test: ConfigureNetworks() -> [application config is nil or empty]
	//test: ConfigureNetworks() -> [application config is nil or empty]
	//test: ConfigureNetworks() -> [file name is empty for case officer: test]
	//test: ConfigureNetworks() -> [agent lookup is nil for case officer: test]

}

func ExampleConfigureNetworks() {
	//var appCfg map[string]string

	appCfg, err := ReadAppConfig(func() ([]byte, error) {
		return readFile(appFileName)
	})
	if err != nil {
		fmt.Printf("test: ReadAppConfig(\"%v\") -> [map:%v] [err:%v]\n", subDir+appFileName, len(appCfg), err)
	}

	errs := ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> [count:%v] [errs:%v]\n", len(errs), errs)

	a := exchange.Agent(module.NamespaceNamePrimary)
	if officer, ok := any(a).(caseofficer.Agent); ok {
		officer.Trace()
	}

	a = exchange.Agent(module.NamespaceNameSecondary)
	if officer, ok := any(a).(caseofficer.Agent); ok {
		officer.Trace()
	}

	fmt.Printf("trace: Operations() -> %v\n", opsAgent.ex.List())

	//Output:
	//test: ConfigureNetworks() -> [count:0] [errs:[]]
	//trace: operative -> test:core:agent/log/access/http
	//trace: operative -> test:resiliency:agent/cache/request/http
	//trace: operative -> test:resiliency:agent/rate-limiting/request/http
	//trace: operative -> test:resiliency:agent/routing/request/http
	//trace: Operations() -> [test:resiliency:agent/caseOfficer/service/traffic/ingress/primary]

}

func ExampleReadAppConfig() {
	cfg, err := ReadAppConfig(func() ([]byte, error) {
		return readFile(appFileName)
	})

	fmt.Printf("test: ReadAppConfig() -> %v [err:%v]\n", cfg, err)

	//Output:
	//test: ReadAppConfig() -> map[test:resiliency:agent/caseOfficer/service/traffic/ingress/primary:network-config-primary.json] [err:<nil>]

}
