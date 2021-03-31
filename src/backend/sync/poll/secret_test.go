package poll

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"

	v1 "k8s.io/api/core/v1"
)

func TestNewSecretPoller(t *testing.T) {
	poller := NewSecretPoller("test-secret", "test-ns", nil)

	if poller == nil {
		t.Fatal("Expected poller not to be nil.")
	}
}

func TestNewSecretPoller_Poll(t *testing.T) {
	var watchEvent *watch.Event
	sName := "test-secret"
	nsName := "test-ns"
	event := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nsName,
			Name:      sName,
		},
	}
	client := fake.NewSimpleClientset(event)
	poller := NewSecretPoller(sName, nsName, client)

	watcher := poller.Poll(1 * time.Second)
	select {
	case ev := <-watcher.ResultChan():
		watchEvent = &ev
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout while waiting for watcher data.")
	}

	if watchEvent == nil {
		t.Fatal("Expected watchEvent not to be nil.")
	}
}
