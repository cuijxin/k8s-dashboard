package sync

import (
	syncApi "github.com/cuijxin/k8s-dashboard/src/backend/sync/api"

	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Implements SynchronizerManager interface.
type synchronizerManager struct {
	client kubernetes.Interface
}

// Secret implements synchronizer manager. See SynchronizerManager interface for more information.
func (s *synchronizerManager) Secret(namespace, name string) syncApi.Synchronizer {
	return &secretSynchronizer{
		namespace:      namespace,
		name:           name,
		client:         s.client,
		actionHandlers: make(map[watch.EventType][]syncApi.ActionHandlerFunction),
	}
}

// NewSynchronizerManager creates new instance of SynchronizerManager.
func NewSynchronizerManager(client kubernetes.Interface) syncApi.SynchronizerManager {
	return &synchronizerManager{client: client}
}
