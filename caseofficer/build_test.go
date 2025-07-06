package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	"github.com/behavioral-ai/core/rest"
	"net/http"
	"reflect"
)

const (
	loggingRole       = "logging"
	authorizationRole = "authorization"
	cacheRole         = "cache"
	rateLimitingRole  = "rate-limiting"
	routingRole       = "routing"
	authorizationName = "Authorization"
	namespaceNameAuth = "test:resiliency:handler/authorization/http"
)

var (
	roles = []string{loggingRole, authorizationRole, cacheRole, rateLimitingRole, routingRole}
)

func authorization(next rest.Exchange) rest.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		auth := r.Header.Get(authorizationName)
		if auth == "" {
			return &http.Response{StatusCode: http.StatusUnauthorized}, nil
		}
		return next(r)
	}
}

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
	//test: buildLink("any:any:link/test/one") -> [<nil>] [err:invalid Namespace kind: link and role: test-role]
	//test: buildLink("any:any:agent/test/one") -> [<nil>] [err:agent is nil for name: any:any:agent/test/one and role: test-role]

}

func ExampleBuildLink() {
	name := "any:any:handler/test/one"
	role := "test-role"
	cfg := make(map[string]string)
	cfg[NameKey] = "any:any:handler/test/one"

	agent := messagingtest.NewAgent("agent\test")
	exchange.RegisterExchangeHandler(name, authorization)

	cfg[NameKey] = name
	t, err := buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", reflect.TypeOf(t), err)

	name = "any:any:agent/test/one"
	cfg[NameKey] = name
	exchange.RegisterConstructor(name, func() messaging.Agent {
		return messagingtest.NewAgent("agent\test2")
	})
	t, err = buildLink(role, cfg, agent)
	fmt.Printf("test: buildLink() -> [%v] [err:%v]\n", reflect.TypeOf(t), err)

	//Output:
	//test: buildLink() -> [func(rest.Exchange) rest.Exchange] [err:<nil>]
	//test: buildLink() -> [*messagingtest.AgentT] [err:<nil>]

}

func ExampleBuildNetwork_Error() {
	name := "*:*:link/test/one"
	officer := messagingtest.NewAgent("*:*:agent/test")
	netCfg := make(map[string]map[string]string)

	chain, errs := buildNetwork(nil, nil, nil)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	chain, errs = buildNetwork(officer, nil, nil)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	chain, errs = buildNetwork(officer, netCfg, roles)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	exchange.RegisterExchangeHandler(name, authorization)

	netCfg[authorizationRole] = map[string]string{}
	chain, errs = buildNetwork(officer, netCfg, roles)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", chain, errs)

	//Output:
	//test: buildNetwork() -> [chain:[]] [agent is nil]
	//test: buildNetwork() -> [chain:[]] [network configuration is nil or empty]
	//test: buildNetwork() -> [chain:[]] [network configuration is nil or empty]
	//test: buildNetwork() -> [chain:[]] [agent or exchange name not found or is empty for role: authorization]

}

func ExampleBuildNetwork() {
	officer := messagingtest.NewAgent("*:*:agent/test")
	netCfg := make(map[string]map[string]string)

	exchange.RegisterExchangeHandler(namespaceNameAuth, authorization)

	netCfg[authorizationRole] = map[string]string{NameKey: namespaceNameAuth}
	chain, errs := buildNetwork(officer, netCfg, roles)
	fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", len(chain), errs)

	//netCfg[routingRole] = map[string]string{NameKey: routing.NamespaceName}
	//chain, errs = buildNetwork(officer, netCfg, roles)
	//fmt.Printf("test: buildNetwork() -> [chain:%v] %v\n", len(chain), errs)

	//Output:
	//test: buildNetwork() -> [chain:1] []

}
