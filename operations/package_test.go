package operations

import (
	"fmt"
	"github.com/appellative-ai/agency/caseofficer"
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/core/std"
	"github.com/appellative-ai/traffic/logger"
	"os"
)

const (
	subDir            = "/operationstest/resource/"
	operatorsFileName = "logging-operators.json"
	endpointFileName  = "endpoint-config.json"
	originFileName    = "origin-config.json"
)

func readFile(fileName string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(dir + subDir + fileName)
}

/*
func readEndpointConfig(read func() ([]byte, error)) ([]map[string]string, error) {
	var cfg []map[string]string

	buf, err := read()
	if err != nil {
		return nil, err //fmt.Printf("test: readFile(\"%v\") -> [bytes:%v] [err:%v]\n", subDir+appFileName, len(buf), err)
	}
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err //fmt.Printf("test: json.Unmarshal() -> [err:%v]\n", err)
	}
	return cfg, nil
}


*/

func ExampleConfigureLogging() {
	exchange.Agent(logger.NamespaceName) //access.Log(nil, "", time.Now().UTC(), 0, "", nil, nil)
	err := ConfigureLogging(func() ([]byte, error) {
		return readFile(operatorsFileName)
	})
	fmt.Printf("test: ConfigureLogging(\"%v\") -> [err:%v]\n", subDir+operatorsFileName, err)

	//Output:
	//test: ConfigureLogging("/operationstest/resource/logging-operators.json") -> [err:<nil>]

}

func ExampleConfigureNetworks_Errors() {
	var appCfg []map[string]string

	errs := ConfigureNetworks(appCfg, nil)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	errs = ConfigureNetworks(nil, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg = []map[string]string{make(map[string]string)}
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	appCfg[0] = make(map[string]string)
	appCfg[0]["endpoint"] = "test"
	errs = ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> %v\n", errs)

	//Output:
	//test: ConfigureNetworks() -> [network configuration read function is nil]
	//test: ConfigureNetworks() -> [endpoint configuration is nil or empty]
	//test: ConfigureNetworks() -> [endpoint name is empty]
	//test: ConfigureNetworks() -> [network file name is empty for endpoint: test]

}

func ExampleConfigureNetworks() {
	appCfg, err := ReadEndpointConfig(func() ([]byte, error) {
		return readFile(endpointFileName)
	})
	if err != nil {
		fmt.Printf("test: ReadEndpointConfig(\"%v\") -> [map:%v] [err:%v]\n", subDir+endpointFileName, len(appCfg), err)
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
	//trace: core:common:agent/caseofficer/request/http/primary -> test:resiliency:agent/cache/request/http
	//trace: core:common:agent/caseofficer/request/http/primary -> test:resiliency:agent/rate-limiting/request/http
	//trace: core:common:agent/caseofficer/request/http/secondary -> test:resiliency:agent/routing/request/http
	//trace: Operations() -> [core:common:agent/caseofficer/request/http/primary core:common:agent/caseofficer/request/http/secondary]

}

/*
func ExampleReadEndpointConfig() {
	cfg, err := ReadEndpointConfig(func() ([]byte, error) {
		return readFile(endpointFileName)
	})

	fmt.Printf("test: ReadEndpointConfig() -> %v [err:%v]\n", cfg, err)

	//Output:
	//test: ReadEndpointConfig() -> map[primary:map[network:network-config-primary.json pattern:/primary] secondary:map[network:network-config-secondary.json pattern:/secondary]] [err:<nil>]

}


*/

var (
	o = std.OriginT{
		Region:      "region",
		Zone:        "zone",
		SubZone:     "sub-zone",
		Host:        "host",
		ServiceName: "service-name",
		InstanceId:  "instance-id",
		Collective:  "collective",
		Domain:      "domain",
	}
)

func ExampleConfigureOrigin() {
	m := map[string]string{
		std.RegionKey:     "us-west1",
		std.ZoneKey:       "oregon",
		std.SubZoneKey:    "portland",
		std.InstanceIdKey: "123456789",
	}
	read := func() ([]byte, error) {
		return readFile(originFileName)
	}
	err := ConfigureOrigin(m, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", subDir+originFileName, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", std.Origin, std.Origin.Host)

	m2 := make(map[string]string)
	err = ConfigureOrigin(m2, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", subDir+originFileName, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", std.Origin, std.Origin.Host)

	m2 = map[string]string{
		std.RegionKey: "us-west1",
		//std.ZoneKey:    "oregon",
		std.SubZoneKey: "portland",
		//std.InstanceIdKey: "123456789",
	}
	err = ConfigureOrigin(m2, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", subDir+originFileName, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", std.Origin, std.Origin.Host)

	//Output:
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:<nil>]
	//test: messaging.SetOrigin() -> google-collective:search:service/us-west1/oregon/portland/google-search#123456789 [host:www.google.com]
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:config map does not contain key: region]
	//test: messaging.SetOrigin() ->  [host:]
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:config map does not contain key: zone]
	//test: messaging.SetOrigin() ->  [host:]

}
