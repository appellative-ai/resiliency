package operations

import (
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
)

const (
	networkFileName = "network-config-primary.json"
)

func ExampleBuildNetworkConfig() {
	cfg, err := buildNetworkConfig(networkFileName, readFile)
	fmt.Printf("test: buildNetworkConfig() -> [%v] [err:%v]\n", cfg, err)

	//Output:
	//test: buildNetworkConfig() -> [cfg:map[authorization:map[name:test:resiliency:link/authorization/http] cache:map[cache-control:no-store, no-cache, max-age=0 fri:22-23 host:localhost:8081 interval:4m mon:8-16 name:test:resiliency:agent/cache/request/http sat:3-8 sun:13-15 thu:0-23 timeout:750ms tue:6-10 wed:12-12] logging:map[name:test:resiliency:link/logging/access] rate-limiting:map[load-size:567 name:test:resiliency:agent/rate-limiting/request/http off-peak-duration:5m peak-duration:750ms rate-burst:12 rate-limit:1234] routing:map[app-host:localhost:8080 cache-host:localhost:8081 interval:4m name:test:resiliency:agent/routing/request/http timeout:2m]]] [err:<nil>]

}

func _ExampleCreateNetworkConfig() {
	var buf []byte
	var err error
	var appCfg []map[string]string

	readNetworkConfig(networkFileName, readFile, &buf, &err)
	fmt.Printf("test: readNetworkConfig(\"%v\") -> [buf:%v] [err:%v]\n", networkFileName, len(buf), err)

	err = json.Unmarshal(buf, &appCfg)
	fmt.Printf("test: json.Unmarshal() -> [err:%v]\n", err)

	cfg := shapeNetworkConfig(appCfg)
	for k, v := range cfg {
		fmt.Printf("test: shapeNetworkConfig() [k:%v] [role:%v] [%v]\n", k, v[RoleKey], v[NameKey])
	}

	//Output:
	//test: readNetworkConfig("network-config-primary.json") -> [buf:995] [err:<nil>]
	//test: json.Unmarshal() -> [err:<nil>]
	//test: shapeNetworkConfig() [k:routing] [role:] [test:resiliency:agent/routing/request/http]
	//test: shapeNetworkConfig() [k:authorization] [role:] [test:resiliency:link/authorization/http]
	//test: shapeNetworkConfig() [k:logging] [role:] [test:resiliency:link/logging/access]
	//test: shapeNetworkConfig() [k:rate-limiting] [role:] [test:resiliency:agent/rate-limiting/request/http]
	//test: shapeNetworkConfig() [k:cache] [role:] [test:resiliency:agent/cache/request/http]

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
