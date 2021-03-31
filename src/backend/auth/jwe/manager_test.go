package jwe

import (
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/sync"
	"reflect"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd/api"
)

func getTokenManager() authApi.TokenManager {
	c := fake.NewSimpleClientset()
	syncManager := sync.NewSynchronizerManager(c)
	holder := NewRSAKeyHolder(syncManager.Secret("", ""))
	return NewJWETokenManager(holder)
}

func areErrorsEqual(err1, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() == err2.Error()) ||
		(err1 == nil && err2 == nil)
}

func TestJewTokenManager_Generate(t *testing.T) {
	cases := []struct {
		info        string
		authInfo    api.AuthInfo
		expectedErr error
	}{
		{
			"Shoulde generate encrypted token",
			api.AuthInfo{Token: "test-token"},
			nil,
		},
	}

	for _, c := range cases {
		tokenManager := getTokenManager()
		token, err := tokenManager.Generate(c.authInfo)

		if !areErrorsEqual(err, c.expectedErr) {
			t.Errorf("Test Case: %s. Expected error to be: %v, but got %v.",
				c.info, c.expectedErr, err)
		}

		if len(token) == 0 {
			t.Errorf("Test Case: %s. Expected token not to be empty.", c.info)
		}
	}
}

func TestJewTokenManager_Decrypt(t *testing.T) {
	cases := []struct {
		info        string
		authInfo    api.AuthInfo
		expected    *api.AuthInfo
		expectedErr error
	}{
		{
			"Should decrypt encrypted token",
			api.AuthInfo{Token: "test-token"},
			&api.AuthInfo{Token: "test-token"},
			nil,
		},
	}

	for _, c := range cases {
		tokenManager := getTokenManager()
		token, _ := tokenManager.Generate(c.authInfo)
		authInfo, err := tokenManager.Decrypt(token)

		if !areErrorsEqual(err, c.expectedErr) {
			t.Errorf("Test Case: %s. Expected error to be: %v, but got %v.",
				c.info, c.expectedErr, err)
		}

		if !reflect.DeepEqual(authInfo, c.expected) {
			t.Errorf("Test Case: %s. Expected: %v, but got %v.", c.info, c.expected,
				authInfo)
		}
	}
}

func TestJweTokenManager_Refresh(t *testing.T) {
	cases := []struct {
		info        string
		authInfo    api.AuthInfo
		shouldSleep bool
		expected    bool
		expectedErr error
	}{
		{
			"Should refresh valid token",
			api.AuthInfo{Token: "test-token"},
			false,
			true,
			nil,
		},
		{
			info:        "Should return error when no token provided",
			authInfo:    api.AuthInfo{},
			shouldSleep: false,
			expected:    false,
			expectedErr: errors.NewInvalid("Can not refresh token. No token provided."),
		},
		{
			info:        "Should return error when token has expired",
			authInfo:    api.AuthInfo{Token: "test-token"},
			shouldSleep: true,
			expected:    false,
			expectedErr: errors.NewTokenExpired(errors.MsgTokenExpiredError),
		},
	}

	for _, c := range cases {
		tokenManager := getTokenManager()
		tokenManager.SetTokenTTL(1)
		token, _ := tokenManager.Generate(c.authInfo)

		if len(c.authInfo.Token) == 0 {
			token = ""
		}

		if c.shouldSleep {
			time.Sleep(2 * time.Second)
		}

		refreshedToken, err := tokenManager.Refresh(token)

		if !areErrorsEqual(err, c.expectedErr) {
			t.Errorf("Test Case: %s. Expected error to be: %v, but got %v.",
				c.info, c.expectedErr, err)
		}

		if (c.expected && len(refreshedToken) == 0) || (!c.expected && len(refreshedToken) > 0) {
			t.Errorf("Test Case: %s. Expected new token to be generated: %t", c.info, c.expected)
		}
	}
}
