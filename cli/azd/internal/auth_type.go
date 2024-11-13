package internal

// AuthType defines different authentication types.
type AuthType string

const (
	AUTH_TYPE_UNSPECIFIED AuthType = "UNSPECIFIED"
	// Username and password, or key based authentication
	AuthType_PASSWORD AuthType = "PASSWORD"
	// Connection string authentication
	AuthType_CONNECTION_STRING AuthType = "CONNECTION_STRING"
	// Microsoft EntraID token credential
	AuthType_MANAGED_IDENTITY AuthType = "MANAGED_IDENTITY"
)
