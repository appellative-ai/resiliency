package endpoint

import "fmt"

func ExampleBuild() {
	err := Build("test", nil)
	fmt.Printf("test: Build() -> [err:%v]\n", err)

	err = Build("test", []any{"test"})
	fmt.Printf("test: Build() -> [err:%v]\n", err)

	//Output:
	//test: Build() -> [err:chain is nil or empty]
	//test: Build() -> [err:agent not found for name: test]

}
