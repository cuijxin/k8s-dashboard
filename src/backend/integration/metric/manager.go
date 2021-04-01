package metric

import (
	"fmt"
	"log"
	"time"

	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
	integrationapi "github.com/cuijxin/k8s-dashboard/src/backend/integration/api"
	metricapi "github.com/cuijxin/k8s-dashboard/src/backend/integration/metric/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration/metric/heapster"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration/metric/sidecar"

	"k8s.io/apimachinery/pkg/util/wait"
)

// MetricManager is responsible for management of all integrated applications related to metrics.
type MetricManager interface {
	// AddClient adds metric client to client list supported by this manager.
	AddClient(metricapi.MetricClient) MetricManager
	// Client returns active Metric client.
	Client() metricapi.MetricClient
	// Enable is responsible for switching active client if given integration application id
	// is found and related application is healthy (we can connect to it).
	Enable(integrationapi.IntegrationID) error
	// EnableWithRetry works similar to enable. It runs in a separate thread and tries to enable integration with given
	// id every 'period' seconds.
	EnableWithRetry(id integrationapi.IntegrationID, period time.Duration)
	// List returns list of available metric related integrations.
	List() []integrationapi.Integration
	// ConfigureSidecar configures and adds sidecar to clients list.
	ConfigureSidecar(host string) MetricManager
	// ConfigureHeapster configures and adds sidecar to clients list.
	ConfigureHeapster(host string) MetricManager
}

// Implements MetricManager interface.
type metricManager struct {
	manager clientapi.ClientManager
	clients map[integrationapi.IntegrationID]metricapi.MetricClient
	active  metricapi.MetricClient
}

// AddClient implements metric manager interface. See MetricManager for more information.
func (m *metricManager) AddClient(client metricapi.MetricClient) MetricManager {
	if client != nil {
		m.clients[client.ID()] = client
	}

	return m
}

// Client implements metric manager interface. See MetricManager for more information.
func (m *metricManager) Client() metricapi.MetricClient {
	return m.active
}

// Enable implements metric manager interface. See MetricManager for more information.
func (m *metricManager) Enable(id integrationapi.IntegrationID) error {
	metricClient, exists := m.clients[id]
	if !exists {
		return fmt.Errorf("No metric client found for integration id: %s", id)
	}

	err := metricClient.HealthCheck()
	if err != nil {
		return fmt.Errorf("Health check failed: %s", err.Error())
	}

	m.active = metricClient
	return nil
}

// EnableWithRetry implements metric manager interface. See MetricManager for more information.
func (m *metricManager) EnableWithRetry(id integrationapi.IntegrationID, period time.Duration) {
	go wait.Forever(func() {
		metricClient, exists := m.clients[id]
		if !exists {
			log.Printf("Metric client with given id %s does not exist.", id)
			return
		}

		err := metricClient.HealthCheck()
		if err != nil {
			m.active = nil
			log.Printf("Metric client health check failed: %s. Retrying in %d seconds.", err, period)
			return
		}

		if m.active == nil {
			log.Printf("Successful request to %s", id)
			m.active = metricClient
		}
	}, period*time.Second)
}

// List implements metric manager interface. See MetricManager for more information.
func (m *metricManager) List() []integrationapi.Integration {
	result := make([]integrationapi.Integration, 0)
	for _, c := range m.clients {
		result = append(result, c.(integrationapi.Integration))
	}

	return result
}

// ConfigureSidecar implements metric manager interface. See MetricManager for more information.
func (m *metricManager) ConfigureSidecar(host string) MetricManager {
	kubeClient := m.manager.InsecureClient()
	metricClient, err := sidecar.CreateSidecarClient(host, kubeClient)
	if err != nil {
		log.Printf("There was an error during sidecar client creation: %s", err.Error())
		return m
	}

	m.clients[metricClient.ID()] = metricClient
	return m
}

// ConfigureHeapster implements metric manager interface. See MetricManager for more information.
func (m *metricManager) ConfigureHeapster(host string) MetricManager {
	kubeClient := m.manager.InsecureClient()
	metricClient, err := heapster.CreateHeapsterClient(host, kubeClient)
	if err != nil {
		log.Printf("There was an error during heapster client creation: %s", err.Error())
		return m
	}

	m.clients[metricClient.ID()] = metricClient
	return m
}

// NewMetricManager creates metric manager.
func NewMetricManager(manager clientapi.ClientManager) MetricManager {
	return &metricManager{
		manager: manager,
		clients: make(map[integrationapi.IntegrationID]metricapi.MetricClient),
	}
}
