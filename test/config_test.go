package test

import "fmt"

const (
	//appConfig     = "file://[cwd]/resource/app-config.json"
	appConfig = "/resource/app-config.json"
	//networkConfig = "file://[cwd]/resource/network-config-primary.json"
	networkConfig = "/resource/network-config-primary.json"
)

func ExampleReadConfig() {
	cfg, err := ReadConfig[map[string]string](appConfig)
	fmt.Printf("test: ReadConfig(\"%v\") -> [map:%v] [%v]\n", appConfig, len(cfg), err)

	cfg2, err2 := ReadConfig[[]map[string]string](networkConfig)
	fmt.Printf("test: readConfig(\"%v\") -> [map:%v] [%v]\n", networkConfig, len(cfg2), err2)

	//Output:
	//test: ReadConfig("/resource/app-config.json") -> [map:2] [<nil>]
	//test: readConfig("/resource/network-config-primary.json") -> [map:4] [<nil>]

}

func ExampleAppConfig() {
	cfg, err := ReadConfig[map[string]string](appConfig)
	if err != nil {
		fmt.Printf("test: ReadConfig(\"%v\") -> [map:%v] [%v]\n", appConfig, len(cfg), err)
	}

	key := "primary"
	app := NewAppConfig(cfg)
	s, ok := app.Name(key)
	fmt.Printf("test: app.Name(\"%v\") -> [%v] [ok:%v]\n", key, s, ok)
	s, ok = app.Path(key)
	fmt.Printf("test: app.Path(\"%v\") -> [%v] [ok:%v]\n", key, s, ok)

	key = "secondary"
	//app := NewAppConfig(cfg)
	s, ok = app.Name(key)
	fmt.Printf("test: app.Name(\"%v\") -> [%v] [ok:%v]\n", key, s, ok)
	s, ok = app.Path(key)
	fmt.Printf("test: app.Path(\"%v\") -> [%v] [ok:%v]\n", key, s, ok)

	//Output:
	//test: app.Name("primary") -> [test:resiliency:agent/caseOfficer/service/traffic/ingress/primary] [ok:true]
	//test: app.Path("primary") -> [network-config-primary.json] [ok:true]
	//test: app.Name("secondary") -> [test:resiliency:agent/caseOfficer/service/traffic/ingress/secondary] [ok:true]
	//test: app.Path("secondary") -> [network-config-secondary.json] [ok:true]

}
