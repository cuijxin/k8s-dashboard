package poll

import (
	"sync"

	"k8s.io/apimachinery/pkg/watch"
)

// Watcher Implements watch.Interface
type Watcher struct {
	eventChan chan watch.Event
	stopped   bool
	sync.Mutex
}

// Stop stops poll watcher and closes event channel.
func (p *Watcher) Stop() {
	p.Lock()
	defer p.Unlock()
	if p.stopped {
		close(p.eventChan)
		p.stopped = true
	}
}

// IsStopped returns whether or not watcher was stopped.
func (p *Watcher) IsStopped() bool {
	return p.stopped
}

// ResultChan returns result channel that user can watch for incoming events.
func (p *Watcher) ResultChan() <-chan watch.Event {
	p.Lock()
	defer p.Unlock()
	return p.eventChan
}

// NewWatcher creates instance of Watcher.
func NewWatcher() *Watcher {
	return &Watcher{
		eventChan: make(chan watch.Event),
		stopped:   false,
	}
}
