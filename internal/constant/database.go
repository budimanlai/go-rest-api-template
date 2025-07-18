package constant

// Database connection constants
const (
	// Connection parameters
	MaxOpenConnections    = 25
	MaxIdleConnections    = 25
	ConnectionMaxLifetime = 5 // minutes
	ConnectionMaxIdleTime = 5 // minutes

	// Database timeouts (in seconds)
	QueryTimeout   = 30
	ConnectTimeout = 10

	// Additional user status (besides the ones in constant.go)
	UserStatusDeleted = "deleted"
	UserStatusBanned  = "banned"

	// Additional defaults
	DefaultCreatedBySystem = "system" // string version
	DefaultUpdatedBySystem = "system" // string version

	// Additional validation constraints
	MaxTokenLength = 255
)
