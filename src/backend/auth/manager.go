package auth

import (
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
)

// Implements AuthManager interface
type authManager struct {
	tokenManager            authApi.TokenManager
	clientManger            clientapi.ClientManager
	authenticationModes     authApi.AuthenticationModes
	authenticationSkippable bool
}
