package endpoint

const (
	OperationsPattern = "/operations"
	RootPattern       = "/"
)

var (
	Operations = NewOperationsEndpoint()
	Root       = NewRootEndpoint()
)
