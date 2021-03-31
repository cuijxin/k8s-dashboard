package poll

import "testing"

func TestNewWatcher(t *testing.T) {
	watcher := NewWatcher()

	if watcher == nil {
		t.Fatal("Expected watcher not to be nil.")
	}
}
