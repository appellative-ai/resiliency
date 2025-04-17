package opstest

import (
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/fmtx"
	"github.com/behavioral-ai/core/iox"
)

const (
	agentsConfigPath  = "file://[cwd]/resource/agents-config.txt"
	originConfigPath  = "file://[cwd]/resource/origin-config.txt"
	loggingConfigPath = "file://[cwd]/resource/logging-operators.json"

	RegionKey     = "region"
	ZoneKey       = "zone"
	SubZoneKey    = "sub-zone"
	HostKey       = "host"
	InstanceIdKey = "instance-id"
)

// Operator - configuration of logging entries
type Operator struct {
	Name  string
	Value string
}

func ExampleAgentsConfig() {
	m, err := iox.ReadMap(agentsConfigPath)
	fmt.Printf("test: iox.ReadMap() -> [err:%v]\n", err)

	timeout := m["timeout"]
	if timeout != "" {
		dur, err1 := fmtx.ParseDuration(timeout)
		fmt.Printf("test: fmtx.ParseDuration(\"%v\") -> [dur:%v] [err:%v]\n", timeout, dur, err1)
	}

	fmt.Printf("test: ReadMap() -> %v\n", m)

	//Output:
	//test: iox.ReadMap() -> [err:<nil>]
	//test: fmtx.ParseDuration("1500ms") -> [dur:1.5s] [err:<nil>]
	//test: ReadMap() -> map[app-host:localhost:8080 cache-host:localhost:8082 timeout:1500ms]

}

func ExampleOriginConfig() {
	m, err := iox.ReadMap(originConfigPath)
	fmt.Printf("test: iox.ReadMap() -> [err:%v]\n", err)

	fmt.Printf("test: Origin() -> [region:%v] [zone:%v] [sub-zone:%v] [host:%v] [instance-id:%v]\n", m[RegionKey], m[ZoneKey], m[SubZoneKey], m[HostKey], m[InstanceIdKey])

	//Output:
	//test: iox.ReadMap() -> [err:<nil>]
	//test: Origin() -> [region:region] [zone:zone] [sub-zone:sub-zone] [host:host-name] [instance-id:instance-id]

}

func ExampleLoggingOperators() {
	buf, err := iox.ReadFile(loggingConfigPath)
	fmt.Printf("test: iox.ReadMap() -> [err:%v]\n", err)

	var ops []Operator
	err = json.Unmarshal(buf, &ops)
	fmt.Printf("test: json.Unmarshal() -> [ops:%v] [err:%v]\n", ops, err)

	//Output:
	//test: iox.ReadMap() -> [err:<nil>]
	//test: json.Unmarshal() -> [ops:[{ %START_TIME%} { %DURATION%} { %TRAFFIC%} {route-name %ROUTE%} { %METHOD%} { %HOST%} { %PATH%} { %STATUS_CODE%} { %TIMEOUT_DURATION%} { %RATE_LIMIT%} { %REDIRECT%}]] [err:<nil>]

}
