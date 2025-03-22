package cache

// Configuration map keys

const (
	CacheHostKey = "cache-host" // Cache service host name
	TimeoutKey   = "timeout"
)

var (
	Agent = New()
)
