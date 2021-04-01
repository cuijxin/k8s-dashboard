package integration

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/integration/api"
)

// IntegrationsGetter is responsible for listing all supported integrations.
type IntegrationsGetter interface {
	// List returns list of all supported integrations.
	List() []api.Integration
}

// List implements integration getter interface. See IntegrationsGetter for more
// information.
func (m *integrationManager) List() []api.Integration {
	result := make([]api.Integration, 0)

	// Append all types of integrations
	result = append(result, m.Metric().List()...)

	return result
}
