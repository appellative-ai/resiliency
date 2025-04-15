package endpoint

const (
	OperationsPattern = "/operations"
	RootPattern       = "/"
	HealthPattern     = "/health"
)

var (
	Operations = NewOperationsEndpoint()
	Root       = NewRootEndpoint()
	Health     = NewHealthEndpoint()
)
