package cache

// Configuration map keys

const (
	HostKey    = "cache-host" // Cache service host name
	TimeoutKey = "timeout"
)

var (
	Agent = New()
)
