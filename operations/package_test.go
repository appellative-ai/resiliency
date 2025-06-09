package operations

import (
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"os"
)

var (
	//originConfig    = "file://[cwd]/operationstest/resource/origin-config.json"
	//operatorsConfig = "file://[cwd]/operationstest/resource/logging-operators.json"

	originPath    = "/operationstest/resource/origin-config.json"
	operatorsPath = "/operationstest/resource/logging-operators.json"

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
	return os.ReadFile(dir + fileName)
}

func ExampleOrigin_Config() {
	m := map[string]string{
		messaging.RegionKey:     "us-west1",
		messaging.ZoneKey:       "oregon",
		messaging.SubZoneKey:    "portland",
		messaging.InstanceIdKey: "123456789",
	}
	read := func() ([]byte, error) {
		return readFile(originPath)
	}
	err := ConfigureOrigin(m, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", originPath, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	m2 := make(map[string]string)
	err = ConfigureOrigin(m2, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", originPath, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	m2 = map[string]string{
		messaging.RegionKey: "us-west1",
		//messaging.ZoneKey:    "oregon",
		messaging.SubZoneKey: "portland",
		//messaging.InstanceIdKey: "123456789",
	}
	err = ConfigureOrigin(m2, read)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", originPath, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	//Output:
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:<nil>]
	//test: messaging.SetOrigin() -> google-collective:search:service/us-west1/oregon/portland/google-search#123456789 [host:www.google.com]
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:config map does not contain key: region]
	//test: messaging.SetOrigin() ->  [host:]
	//test: ConfigOrigin("/operationstest/resource/origin-config.json") -> [err:config map does not contain key: zone]
	//test: messaging.SetOrigin() ->  [host:]

}

func _ExampleOrigin_Map() {
	m := map[string]string{
		messaging.RegionKey:      "region",
		messaging.ZoneKey:        "zone",
		messaging.SubZoneKey:     "sbu-zone",
		messaging.HostKey:        "host",
		messaging.ServiceNameKey: "service-name",
		messaging.InstanceIdKey:  "instance-id",
		messaging.CollectiveKey:  "collective",
		messaging.DomainKey:      "domain",
	}

	buf, err := json.Marshal(&m)
	fmt.Printf("test: Marshal() -> [err:%v]\n", err)
	fmt.Printf("test: Marshal() -> %v\n", string(buf))

	var m2 map[string]string

	err = json.Unmarshal(buf, &m2)
	fmt.Printf("test: Unmarshal() -> [err:%v]\n", err)
	fmt.Printf("test: Unmarshal() -> %v\n", m2)

	//Output:
	//fail

}

func ExampleLogging_Config() {
	err := ConfigureLogging(func() ([]byte, error) {
		return readFile(operatorsPath)
	})
	fmt.Printf("test: ConfigureLogging(\"%v\") -> [err:%v]\n", operatorsPath, err)

	//Output:
	//test: ConfigureLogging("/operationstest/resource/logging-operators.json") -> [err:<nil>]

}

//buf, err := iox.ReadFile(originConfig)
//fmt.Printf("test: iox.ReadFile(\"%v\") -> [err:%v]\n", originConfig, err)
//err = json.Unmarshal(buf, &m)
//fmt.Printf("test: Unmarshal() -> [err:%v]\n", err)
//fmt.Printf("test: Unmarshal() -> %v\n", m)
