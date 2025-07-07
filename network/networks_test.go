package network

import (
	"fmt"
	"os"
)

const (
	networkFileName = "network-config-primary.json"
	subDir          = "/networktest/resource/"
)

func readFile(fileName string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(dir + subDir + fileName)
}

func ExampleBuildNetworkConfig() {
	cfg, err := BuildConfig(networkFileName, readFile)
	fmt.Printf("test: buildNetworkConfig() -> [%v] [err:%v]\n", cfg, err)

	//Output:
	//test: buildNetworkConfig() -> [map[authorization:map[name:test:resiliency:handler/authorization/http] cache:map[cache-control:no-store, no-cache, max-age=0 fri:22-23 host:localhost:8081 interval:4m mon:8-16 name:test:resiliency:agent/cache/request/http sat:3-8 sun:13-15 thu:0-23 timeout:750ms tue:6-10 wed:12-12] logging:map[name:test:core:agent/log/access/http] rate-limiting:map[assignment:local load-size:567 name:test:resiliency:agent/rate-limiting/request/http off-peak-duration:5m peak-duration:750ms rate-burst:12 rate-limit:1234] routing:map[app-host:localhost:8080 assignment:global cache-host:localhost:8081 interval:4m name:test:resiliency:agent/routing/request/http timeout:2m]]] [err:<nil>]

}

/*
func ExampleValidateOfficer() {
	agent, err := validateOfficerType("")
	fmt.Printf("test: validateOfficerType() -> [agent:%v] [err:%v]\n", agent, err)

	agent, err = validateOfficerType("test-case-officer")
	fmt.Printf("test: validateOfficerType() -> [agent:%v] [err:%v]\n", agent, err)

	name := "test-agent"
	a := messagingtest.NewAgent(name)
	exchange.Register(a)
	agent, err = validateOfficerType(name)
	fmt.Printf("test: validateOfficerType() -> [agent:%v] [err:%v]\n", agent, err)

	agent, err = validateOfficerType(caseofficer.NamespaceNamePrimary)
	fmt.Printf("test: validateOfficerType() -> [agent:%v] [err:%v]\n", agent != nil, err)

	//Output:
	//test: validateOfficerType() -> [agent:<nil>] [err:case officer name is empty]
	//test: validateOfficerType() -> [agent:<nil>] [err:agent lookup is nil for case officer: test-case-officer]
	//test: validateOfficerType() -> [agent:<nil>] [err:agent is not of type caseofficer.Agent for case officer: test-agent]
	//test: validateOfficerType() -> [agent:true] [err:<nil>]

}


*/
