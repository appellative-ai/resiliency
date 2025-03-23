package opstest

import (
	"fmt"
	"github.com/behavioral-ai/core/fmtx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/resiliency/common"
)

const (
	serviceConfigTxt = "file://[cwd]/operations-config.txt"
)

func ExampleConfig() {
	m, err := iox.ReadMap(serviceConfigTxt)
	fmt.Printf("test: ReadMap() -> [err:%v]\n", err)

	timeout := m[common.TimeoutKey]
	if timeout != "" {
		dur, err1 := fmtx.ParseDuration(timeout)
		fmt.Printf("test: fmtx.ParseDuration(\"%v\") -> [dur:%v] [err:%v]\n", timeout, dur, err1)
	}

	fmt.Printf("test: ReadMap() -> %v\n", m)

	//Output:
	//test: ReadMap() -> [err:<nil>]
	//test: fmtx.ParseDuration("1500ms") -> [dur:1.5s] [err:<nil>]
	//test: ReadMap() -> map[app-host:localhost:8080 cache-host:localhost:8082 host:host-name instance-id:instance-id region:region sub-zone:sub-zone timeout:1500ms zone:zone]

}
