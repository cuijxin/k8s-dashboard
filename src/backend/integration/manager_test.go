package integration

import (
	"strings"
	"testing"

	client "github.com/cuijxin/k8s-dashboard/src/backend/client"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration/api"
)

func areErrorsEqual(err1, err2 error) bool {
	return (err1 != nil && err2 != nil && normalize(err1.Error()) == normalize(err2.Error())) ||
		(err1 == nil && err2 == nil)
}

// Removes all quote signs that might have been added to the message.
// Might depend on dependencies version how they are constructed.
func normalize(msg string) string {
	return strings.Replace(msg, "\"", "", -1)
}

func TestNewIntegrationManager(t *testing.T) {
	iManager := NewIntegrationManager(nil)
	if iManager == nil {
		t.Error("Failed to create integration manager.")
	}
}

func TestIntegrationManager_GetState(t *testing.T) {
	cases := []struct {
		info          string
		apiServerHost string
		heapsterHost  string
		expected      *api.IntegrationState
		expectedErr   error
	}{
		{
			"Server provided and using in-cluster heapster",
			"http://127.0.0.1:8080", "", &api.IntegrationState{
				Connected: false,
				Error:     errors.NewInvalid("Get http://127.0.0.1:8080/api/v1/namespaces/kube-system/services/heapster/proxy/healthz: dial tcp 127.0.0.1:8080: connect: connection refused"),
			}, nil,
		},
		{
			"Server provided and using external heapster",
			"http://127.0.0.1:8080", "http://127.0.0.1:8081", &api.IntegrationState{
				Connected: false,
				Error:     errors.NewInvalid("Get http://127.0.0.1:8081/healthz: dial tcp 127.0.0.1:8081: connect: connection refused"),
			}, nil,
		},
	}

	for _, c := range cases {
		cManager := client.NewClientManager("", c.apiServerHost)
		iManager := NewIntegrationManager(cManager)
		iManager.Metric().ConfigureHeapster(c.heapsterHost)

		state, err := iManager.GetState(api.HeapsterIntegrationID)
		if !areErrorsEqual(err, c.expectedErr) {
			t.Errorf("Test Case: %s. Expected error to be: %v, but got %v.",
				c.info, c.expectedErr, err)
		}

		// Time is irrelevant so we don't need to check it
		if c.expectedErr == nil && (!areErrorsEqual(state.Error, c.expected.Error)) {
			t.Errorf("Test Case: %s. Expected state to be: %v, but got %v.",
				c.info, c.expected.Error, state.Error)
		} else if state.Connected != c.expected.Connected {
			t.Errorf("Test Case: %s. Could not connect to API server.",
				c.info)
		}
	}
}

func TestIntegrationManager_Metric(t *testing.T) {
	metricManager := NewIntegrationManager(nil).Metric()
	if metricManager == nil {
		t.Error("Failed to get metric manager.")
	}
}
