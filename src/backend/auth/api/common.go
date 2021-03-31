package api

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/args"
	"strings"
)

// ToAuthenticationModes transform array of authentication mode strings to valid
// AuthenticationModes type.
func ToAuthenticationModes(modes []string) AuthenticationModes {
	result := AuthenticationModes{}
	modesMap := map[string]bool{}

	for _, mode := range []AuthenticationMode{Token, Basic} {
		modesMap[mode.String()] = true
	}

	for _, mode := range modes {
		if _, exists := modesMap[mode]; exists {
			result.Add(AuthenticationMode(mode))
		}
	}

	return result
}

// List of protected resources that should be filtered out from dashboard UI.
var protectedResources = []ProtectedResource{
	{EncryptionKeyHolderName, args.Holder.GetNamespace()},
	{CertificateHolderSecretName, args.Holder.GetNamespace()},
}

// ShouldRejectRequest returns true if url contains name and namespace of resource
// that should be filtered out from dashboard.
func ShouldRejectRequest(url string) bool {
	for _, protectedResource := range protectedResources {
		if strings.Contains(url, protectedResource.ResourceName) && strings.Contains(url, protectedResource.ResourceNamespace) {
			return true
		}
	}
	return false
}
