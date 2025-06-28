package auth

// Authenticator is the interface that wraps the basic Authenticate method.
type Authenticator interface {
	Authenticate(username, password string) (string, error)
}
