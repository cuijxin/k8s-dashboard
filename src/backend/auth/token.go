package auth

import (
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"

	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements Authenticator interface
type tokenAuthenticator struct {
	token string
}

// GetAuthInfo implements Authenticator interface. See Authenticator for more
// information.
func (t tokenAuthenticator) GetAuthInfo() (api.AuthInfo, error) {
	return api.AuthInfo{
		Token: t.token,
	}, nil
}

// NewTokenAuthenticator returns Authenticator based on LoginSpec.
func NewTokenAuthenticator(spec *authApi.LoginSpec) authApi.Authenticator {
	return &tokenAuthenticator{
		token: spec.Token,
	}
}
