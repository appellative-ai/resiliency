package operations

import (
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
)

var (
	originConfig = "file://[cwd]/operationstest/resource/origin-config.json"

	o = messaging.OriginT{
		Name:        "name",
		Region:      "region",
		Zone:        "zone",
		SubZone:     "sbu-zone",
		Host:        "host",
		ServiceName: "service-name",
		InstanceId:  "instance-id",
		Collective:  "collective",
		Domain:      "domain",
	}
)

func ExampleOrigin_Config() {
	m := map[string]string{
		messaging.RegionKey:     "us-west1",
		messaging.ZoneKey:       "oregon",
		messaging.SubZoneKey:    "portland",
		messaging.InstanceIdKey: "123456789",
	}

	err := ConfigureOrigin(originConfig, m)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", originConfig, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	m2 := make(map[string]string)
	err = ConfigureOrigin(originConfig, m2)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", originConfig, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	m2 = map[string]string{
		messaging.RegionKey: "us-west1",
		//messaging.ZoneKey:    "oregon",
		messaging.SubZoneKey: "portland",
		//messaging.InstanceIdKey: "123456789",
	}
	err = ConfigureOrigin(originConfig, m2)
	fmt.Printf("test: ConfigOrigin(\"%v\") -> [err:%v]\n", originConfig, err)
	fmt.Printf("test: messaging.SetOrigin() -> %v [host:%v]\n", messaging.Origin, messaging.Origin.Host)

	//Output:
	//test: ConfigOrigin("file://[cwd]/operationstest/resource/origin-config.json") -> [err:<nil>]
	//test: messaging.SetOrigin() -> google-collective:search:service/us-west1/oregon/portland/google-search#123456789 [host:www.google.com]
	//test: ConfigOrigin("file://[cwd]/operationstest/resource/origin-config.json") -> [err:config map does not contain key: region]
	//test: messaging.SetOrigin() ->  [host:]
	//test: ConfigOrigin("file://[cwd]/operationstest/resource/origin-config.json") -> [err:config map does not contain key: zone]
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

//buf, err := iox.ReadFile(originConfig)
//fmt.Printf("test: iox.ReadFile(\"%v\") -> [err:%v]\n", originConfig, err)
//err = json.Unmarshal(buf, &m)
//fmt.Printf("test: Unmarshal() -> [err:%v]\n", err)
//fmt.Printf("test: Unmarshal() -> %v\n", m)
