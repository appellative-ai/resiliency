package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	link "github.com/behavioral-ai/resiliency/link"
	"reflect"
)

func ExampleBuildLink_Error() {
	role := "test-role"
	cfg := make(map[string]string)
	agent := messagingtest.NewAgent("agent\test")

	t, err := buildLink(nil, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", t, err)

	t, err = buildLink(agent, nil, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", t, err)

	t, err = buildLink(agent, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", t, err)

	cfg[NameKey] = "any:any:link/test/one"
	t, err = buildLink(agent, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", t, err)

	cfg[NameKey] = "any:any:agent/test/one"
	t, err = buildLink(agent, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", t, err)

	cfg[NameKey] = "any:any:event/test/one"
	t, err = buildLink(agent, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", t, err)

	//Output:
	//test: buildLink() -> [<nil>] [err:agent or configuration map is nil]
	//test: buildLink() -> [<nil>] [err:agent or configuration map is nil]
	//test: buildLink() -> [<nil>] [err:agent or exchange name not found or is empty for role: test-role]
	//test: buildLink() -> [<nil>] [err:exchange link is nil for name: any:any:link/test/one and role: test-role]
	//test: buildLink() -> [<nil>] [err:agent is nil for name: any:any:agent/test/one and role: test-role]
	//test: buildLink() -> [<nil>] [err:invalid Namespace kind: event and role: test-role]

}

func ExampleBuildLink() {
	name := "any:any:link/test/one"
	role := "test-role"
	cfg := make(map[string]string)
	cfg[NameKey] = "any:any:link/test/one"

	agent := messagingtest.NewAgent("agent\test")
	repository.RegisterExchangeLink(name, link.Logger)

	cfg[NameKey] = name
	t, err := buildLink(agent, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", reflect.TypeOf(t), err)

	name = "any:any:agent/test/one"
	cfg[NameKey] = name
	repository.RegisterConstructor(name, func() messaging.Agent {
		return messagingtest.NewAgent("agent\test2")
	})
	t, err = buildLink(agent, cfg, role)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", reflect.TypeOf(t), err)

	//Output:
	//test: buildLink() -> [rest.ExchangeLink] [err:<nil>]
	//test: buildLink() -> [*messagingtest.agentT] [err:<nil>]

}
