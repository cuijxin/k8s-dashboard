package poll

import (
	"context"
	"time"

	syncapi "github.com/cuijxin/k8s-dashboard/src/backend/sync/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretPoller implements Poller interface. See Poller for more information.
type SecretPoller struct {
	name      string
	namespace string
	client    kubernetes.Interface
	watcher   *Watcher
}

// Poll new secret every 'interval' time and send it to watcher channel. See Poller
// for more information.
func (s *SecretPoller) Poll(interval time.Duration) watch.Interface {
	stopCh := make(chan struct{})

	go wait.Until(func() {
		if s.watcher.IsStopped() {
			close(stopCh)
			return
		}

		s.watcher.eventChan <- s.getSecretEvent()
	}, interval, stopCh)

	return s.watcher
}

// Gets secret from API server and transforms it to watch.Event object.
func (s *SecretPoller) getSecretEvent() (event watch.Event) {
	secret, err := s.client.CoreV1().Secrets(s.namespace).Get(context.TODO(), s.name, metav1.GetOptions{})
	event = watch.Event{
		Object: secret,
		Type:   watch.Added,
	}

	if err != nil {
		event.Type = watch.Error
	}

	// In case it was never created we can still mark it as deleted and let secret
	// be recreated.
	if errors.IsNotFoundError(err) {
		event.Type = watch.Deleted
	}

	return
}

// NewSecretPoller returns instance of Poller interface.
func NewSecretPoller(name, namespace string, client kubernetes.Interface) syncapi.Poller {
	return &SecretPoller{
		name:      name,
		namespace: namespace,
		client:    client,
		watcher:   NewWatcher(),
	}
}
