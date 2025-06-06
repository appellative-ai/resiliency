package operations

import (
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
)

var (
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

func ExampleConfigureOrigin() {
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
