package operations

import (
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"os"
)

var (
	subDir            = "/operationstest/resource/"
	originFileName    = "origin-config.json"
	operatorsFileName = "logging-operators.json"
	appFileName       = "app-config.json"

	o = messaging.OriginT{
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

func readFile(fileName string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(dir + subDir + fileName)
}

func ExampleOrigin_Config() {
	m := map[string]string{
		messaging.RegionKey:     "us-west1",
		messaging.ZoneKey:       "oregon",
		messaging.SubZoneKey:    "portland",
		messaging.InstanceIdKey: "123456789",
	}
	read := func() ([]byte, error) {
		return readFile(originFileName)
	}
	err := ConfigureOrigin(m, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", subDir+originFileName, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	m2 := make(map[string]string)
	err = ConfigureOrigin(m2, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", subDir+originFileName, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	m2 = map[string]string{
		messaging.RegionKey: "us-west1",
		//messaging.ZoneKey:    "oregon",
		messaging.SubZoneKey: "portland",
		//messaging.InstanceIdKey: "123456789",
	}
	err = ConfigureOrigin(m2, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", subDir+originFileName, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	//Output:
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:<nil>]
	//test: messaging.SetOrigin() -> google-collective:search:service/us-west1/oregon/portland/google-search#123456789 [host:www.google.com]
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:config map does not contain key: region]
	//test: messaging.SetOrigin() ->  [host:]
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:config map does not contain key: zone]
	//test: messaging.SetOrigin() ->  [host:]

}

func ExampleLogging_Config() {
	err := ConfigureLogging(func() ([]byte, error) {
		return readFile(operatorsFileName)
	})
	fmt.Printf("test: ConfigureLogging(\"%v\") -> [err:%v]\n", subDir+operatorsFileName, err)

	//Output:
	//test: ConfigureLogging("/operationstest/resource/logging-operators.json") -> [err:<nil>]

}

func ExampleApp_Config() {
	var m map[string]string

	buf, err := readFile(appFileName)
	fmt.Printf("test: readFile(\"%v\") -> [bytes:%v] [err:%v]\n", subDir+appFileName, len(buf), err)

	err = json.Unmarshal(buf, &m)
	fmt.Printf("test: json.Unmarshal() -> %v [err:%v]\n", m, err)

	//Output:
	//test: readFile("/operationstest/resource/app-config.json") -> [bytes:213] [err:<nil>]
	//test: json.Unmarshal() -> map[test:resiliency:agent/caseOfficer/service/traffic/ingress/primary:network-config-primary.json test:resiliency:agent/caseOfficer/service/traffic/ingress/secondary:network-config-secondary.json] [err:<nil>]

}

func _ExampleNetwork_Config() {
	var m []map[string]string

	buf, err := readFile(networkFileName)
	fmt.Printf("test: readFile(\"%v\") -> [bytes:%v] [err:%v]\n", subDir+networkFileName, len(buf), err)

	err = json.Unmarshal(buf, &m)
	fmt.Printf("test: json.Unmarshal() -> %v [err:%v]\n", m, err)

	//Output:
	//test: readFile("/operationstest/resource/network-config-primary.json") -> [bytes:1306] [err:<nil>]
	//test: json.Unmarshal() -> [map[load-size:567 name:test:resiliency:agent/rate-limiting/request/http off-peak-duration:5m peak-duration:750ms rate-burst:12 rate-limit:1234 role:rate-limiting] map[agent:test:resiliency:agent/redirect/request/http interval:4m new-path:/resource/v2 original-path:resource/v1 percentile-threshold:99/1500ms rate-burst:12 rate-limit:1234 role:redirect status-code-threshold:10] map[cache-control:no-store, no-cache, max-age=0 fri:22-23 host:www.google.com interval:4m mon:8-16 name:test:resiliency:agent/cache/request/http role:cache sat:3-8 sun:13-15 thu:0-23 timeout:750ms tue:6-10 wed:12-12] map[app-host:localhost:8082 log:true name:test:resiliency:agent/routing/request/http role:routing route-name:test-route timeout:2m] map[name:test:resiliency:link/authorization/http role:authorization] map[name:test:resiliency:link/logging/access role:logging]] [err:<nil>]

}

func ExampleConfigureNetworks() {
	var appCfg map[string]string

	buf, err := readFile(appFileName)
	fmt.Printf("test: readFile(\"%v\") -> [bytes:%v] [err:%v]\n", subDir+appFileName, len(buf), err)

	err = json.Unmarshal(buf, &appCfg)
	fmt.Printf("test: json.Unmarshal() -> [err:%v]\n", err)

	errs := ConfigureNetworks(appCfg, readFile)
	fmt.Printf("test: ConfigureNetworks() -> [errs:%v]\n", errs)

	//Output:
	//fail
}
