package integration

import (
	"fmt"

	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration/metric"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IntegrationManager is responsible for management of all integration applications.
type IntegrationManager interface {
	// IntegrationsGetter is responsible for listing all supported integrations.
	IntegrationsGetter
	// GetState returns state of integration based on its' id.
	GetState(id api.IntegrationID) (*api.IntegrationState, error)
	// Metric returns metric manager that is responsible management of metric integrations.
	Metric() metric.MetricManager
}

// Implements IntegrationManager interface
type integrationManager struct {
	metric metric.MetricManager
}

// Metric implements integration manager interface. See IntegrationManager for more
// information.
func (m *integrationManager) Metric() metric.MetricManager {
	return m.metric
}

// GetState implements integration manager interface. See IntegrationManager for more
// information.
func (m *integrationManager) GetState(id api.IntegrationID) (*api.IntegrationState, error) {
	for _, i := range m.List() {
		if i.ID() == id {
			return m.getState(i), nil
		}
	}
	return nil, fmt.Errorf("Integration with given id %s does not exist", id)
}

// Checks and returns state of the provided integration application.
func (m *integrationManager) getState(integration api.Integration) *api.IntegrationState {
	result := &api.IntegrationState{
		Error: integration.HealthCheck(),
	}

	result.Connected = result.Error == nil
	result.LastChecked = v1.Now()

	return result
}

// NewIntegrationManager creates integration manager.
func NewIntegrationManager(manager clientapi.ClientManger) IntegrationManager {
	return &integrationManager{
		metric: metric.NewMetricManager(manager),
	}
}
