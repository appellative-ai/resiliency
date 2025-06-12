package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	link "github.com/behavioral-ai/resiliency/link"
	"github.com/behavioral-ai/traffic/routing"
	"reflect"
)

func ExampleBuildLink_Error() {
	//name := "any:any:aspect/test/one"
	role := "test-role"
	cfg := make(map[string]string)
	agent := messagingtest.NewAgent("agent\test")

	t, err := buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink(\"%v\") -> [%v] [err:%v]\n", cfg[NameKey], t, err)

	cfg[NameKey] = "any:any:aspect/test/one"
	t, err = buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink(\"%v\") -> [%v] [err:%v]\n", cfg[NameKey], t, err)

	cfg[NameKey] = "any:any:link/test/one"
	t, err = buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink(\"%v\") -> [%v] [err:%v]\n", cfg[NameKey], t, err)

	cfg[NameKey] = "any:any:agent/test/one"
	t, err = buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink(\"%v\") -> [%v] [err:%v]\n", cfg[NameKey], t, err)

	//Output:
	//test: buildLink("") -> [<nil>] [err:agent or exchange name not found or is empty for role: test-role]
	//test: buildLink("any:any:aspect/test/one") -> [<nil>] [err:invalid Namespace kind: aspect and role: test-role]
	//test: buildLink("any:any:link/test/one") -> [<nil>] [err:exchange link is nil for name: any:any:link/test/one and role: test-role]
	//test: buildLink("any:any:agent/test/one") -> [<nil>] [err:agent is nil for name: any:any:agent/test/one and role: test-role]

}

func ExampleBuildLink() {
	name := "any:any:link/test/one"
	role := "test-role"
	cfg := make(map[string]string)
	cfg[NameKey] = "any:any:link/test/one"

	agent := messagingtest.NewAgent("agent\test")
	repository.RegisterExchangeLink(name, link.Authorization)

	cfg[NameKey] = name
	t, err := buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", reflect.TypeOf(t), err)

	name = "any:any:agent/test/one"
	cfg[NameKey] = name
	repository.RegisterConstructor(name, func() messaging.Agent {
		return messagingtest.NewAgent("agent\test2")
	})
	t, err = buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", reflect.TypeOf(t), err)

	//Output:
	//test: buildLink() -> [rest.ExchangeLink] [err:<nil>]
	//test: buildLink() -> [*messagingtest.AgentT] [err:<nil>]

}

func ExampleBuildNetwork_Error() {
	name := "*:*:link/test/one"
	officer := messagingtest.NewAgent("*:*:agent/test")
	netCfg := make(map[string]map[string]string)

	chain, errs := buildNetwork(nil, nil)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	chain, errs = buildNetwork(officer, nil)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	chain, errs = buildNetwork(officer, netCfg)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	repository.RegisterExchangeLink(name, link.Authorization)

	netCfg[AuthorizationRole] = map[string]string{}
	chain, errs = buildNetwork(officer, netCfg)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	//Output:
	//test: buildNetwork() -> [chain:[]] [error: agent is nil]
	//test: buildNetwork() -> [chain:[]] [error: network configuration is nil or empty]
	//test: buildNetwork() -> [chain:[]] [error: network configuration is nil or empty]
	//test: buildNetwork() -> [chain:[]] [agent or exchange name not found or is empty for role: authorization]

}

func ExampleBuildNetwork() {
	name := "*:*:link/test/one"
	officer := messagingtest.NewAgent("*:*:agent/test")
	netCfg := make(map[string]map[string]string)

	repository.RegisterExchangeLink(name, link.Authorization)

	netCfg[AuthorizationRole] = map[string]string{NameKey: link.NamespaceNameAuth}
	chain, errs := buildNetwork(officer, netCfg)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", len(chain), errs)

	netCfg[RoutingRole] = map[string]string{NameKey: routing.NamespaceName}
	chain, errs = buildNetwork(officer, netCfg)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", len(chain), errs)

	//Output:
	//test: buildNetwork() -> [chain:1] [error: no routing agent was configured]
	//test: buildNetwork() -> [chain:2] []
	
}
