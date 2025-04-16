package endpoint

const (
	OperationsPattern = "/operations"
	RootPattern       = "/"
	HealthPattern     = "/health"
)

var (
	Operations = newOperationsEndpoint()
	Root       = newRootEndpoint()
	Health     = newHealthEndpoint()
)
