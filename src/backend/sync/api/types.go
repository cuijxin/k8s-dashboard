package api

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

// ActionHandlerFunction is a callback function that can be registered on a watch event.
type ActionHandlerFunction func(runtime.Object)

// Synchronizer is used to watch over a kubernetes resource changes in real time.
// It can be used to i.e. Synchronizer encryption key data between multiple dashboard
// replicas.
type Synchronizer interface {
	// Name returns unique name of created synchronizer.
	Name() string
	// Start synchronizer in a separate goruntine. Should not block thread that calls it.
	Start()
	// Error returns error channel. Any error that happens during running synchronizer
	// will be send to this channel.
	Error() chan error
	// Create given runtime object matching synchronized object details (specially
	// type, name, namespace).
	Create(runtime.Object) error
	// Returns local copy of synchronized object or nil in case object has not yet
	// been created or running goroutine did not yet synced it from server.
	Get() runtime.Object
	// Update synchronized object with given object.
	Update(runtime.Object) error
	// Delete synchronized object.
	Delete() error
	// Force synchronous refresh of local object with object got from kubernetes.
	Refresh()
	// RegisterActionHandler registers callback functions on given event types.
	// They are automatically called by watcher.
	RegisterActionHandler(ActionHandlerFunction, ...watch.EventType)
	// SetPoller allows to set custom poller to synchronize objects.
	SetPoller(poller Poller)
}

// SynchronizerManager interface is responsible for creating specific synchronizers.
type SynchronizerManager interface {
	// Secret created single secret synchronizer based on name and namespace information.
	Secret(namespace, name string) Synchronizer
}

// Poller interface is responsible for periodically polling specific resource.
type Poller interface {
	// Poll polls specific resource every 'interval' time. Watch interface is
	// returned in order to use it in the same way as regular watch on resource.
	Poll(interval time.Duration) watch.Interface
}
