package auth

import (
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements AuthManager interface
type authManager struct {
	tokenManager            authApi.TokenManager
	clientManager           clientapi.ClientManager
	authenticationModes     authApi.AuthenticationModes
	authenticationSkippable bool
}

// Login implements auth manager. See AuthManager interface for more information.
func (self authManager) Login(spec *authApi.LoginSpec) (*authApi.AuthResponse, error) {

}

// Refresh implements auth manager. See AuthManager interface for more information.
func (self authManager) Refresh(jweToken string) (string, error) {
	return self.tokenManager.Refresh(jweToken)
}

func (self authManager) AuthenticationModes() []authApi.AuthenticationMode {
	return self.authenticationModes.Array()
}

// Returns authenticator based on provided LoginSpec.
func (self authManager) getAuthenticator(spec *authApi.LoginSpec) (authApi.Authenticator, error) {
	if len(self.authenticationModes) == 0 {
		return nil, errors.NewInvalid("All authentication options disabled. Check --authentication-modes argument for more information.")
	}

	switch {
	case len(spec.Token) > 0 && self.authenticationModes.IsEnabled(authApi.Token):
		return NewTokenAuthenticator(spec), nil
	case len(spec.Username) > 0 && len(spec.Password) > 0 && self.authenticationModes.IsEnabled(authApi.Basic):
		return NewBasicAuthenticator(spec), nil
	case len(spec.KubeConfig) > 0:
		return NewKubeConfigAuthenticator(spec, self.authenticationModes), nil
	}

	return nil, errors.NewInvalid("Not enough data to create authenticator.")
}

// Checks if user data extracted from provided AuthInfo structure is valid and user is correctly authenticated
// by K8S apiserver.
func (self authManager) healthCheck(authInfo api.AuthInfo) error {
	return self.clientManager.HasAccess(authInfo)
}

// NewAuthManager creates auth manager.
func NewAuthManager(clientManager clientapi.ClientManager, tokenManager authApi.TokenManager,
	authenticationModes authApi.AuthenticationModes, authenticationSkippable bool) authApi.AuthManager {
	return &authManager{
		tokenManager:            tokenManager,
		clientManager:           clientManager,
		authenticationModes:     authenticationModes,
		authenticationSkippable: authenticationSkippable,
	}
}
