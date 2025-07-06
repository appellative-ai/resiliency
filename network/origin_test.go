package network

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
)

var (
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
	//subDir            = "/networktest/resource/"
	originFileName = "origin-config.json"
)

func ExampleConfigureOrigin() {
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
	//test: ConfigOrigin("/networktest/resource/origin-config.json") -> [err:<nil>]
	//test: messaging.SetOrigin() -> google-collective:search:service/us-west1/oregon/portland/google-search#123456789 [host:www.google.com]
	//test: ConfigOrigin("/networktest/resource/origin-config.json") -> [err:config map does not contain key: region]
	//test: messaging.SetOrigin() ->  [host:]
	//test: ConfigOrigin("/networktest/resource/origin-config.json") -> [err:config map does not contain key: zone]
	//test: messaging.SetOrigin() ->  [host:]

}
