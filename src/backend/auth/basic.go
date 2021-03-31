package auth

import (
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"

	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements Authenticator interface
type basicAuthenticator struct {
	username string
	password string
}

// GetAuthInfo implements Authenticator interface. See Authenticator for more
// information.
func (b *basicAuthenticator) GetAuthInfo() (api.AuthInfo, error) {
	return api.AuthInfo{
		Username: b.username,
		Password: b.password,
	}, nil
}

// NewBasicAuthenticator returns Authenticator based on LoginSpec.
func NewBasicAuthenticator(spec *authApi.LoginSpec) authApi.Authenticator {
	return &basicAuthenticator{
		username: spec.Username,
		password: spec.Password,
	}
}
