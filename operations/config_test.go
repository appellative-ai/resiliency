package operations

import (
	"fmt"
	"github.com/behavioral-ai/core/fmtx"
	"github.com/behavioral-ai/core/iox"
)

const (
	opsConfigTxt = "file://[cwd]/resource/operations-config.txt"
)

func ExampleConfigure() {
	m, err := iox.ReadMap(opsConfigTxt)
	dur, err1 := fmtx.ParseDuration(m[TimeoutKey])
	fmt.Printf("test: ParseDuration() -> [orig:%v] [parsed:%v][err:%v]\n", m[TimeoutKey], dur, err1)

	fmt.Printf("test: ReadMap() -> %v [err:%v] \n", err, m)

	//Output:
	//test: ParseDuration() -> [orig:1500ms] [parsed:1.5s][err:<nil>]
	//test: ReadMap() -> <nil> [err:map[app-host:localhost:8080 cache-host:localhost:8082 host:host-name instance-id:instance-id region:region sub-zone:sub-zone timeout:1500ms zone:zone]]

}
