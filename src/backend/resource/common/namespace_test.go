package common

import "testing"

func TestToRequestParam(t *testing.T) {
	nsQ := NewSameNamespaceQuery("foo")
	if nsQ.ToRequestParam() != "foo" {
		t.Errorf("Expected %s to be foo", nsQ.ToRequestParam())
	}

	nsQ = NewNamespaceQuery([]string{"foo", "bar"})
	if nsQ.ToRequestParam() != "" {
		t.Errorf("Expected %s to be ''", nsQ.ToRequestParam())
	}

	nsQ = NewNamespaceQuery([]string{})
	if nsQ.ToRequestParam() != "" {
		t.Errorf("Expected %s to be ''", nsQ.ToRequestParam())
	}

	nsQ = NewNamespaceQuery(nil)
	if nsQ.ToRequestParam() != "" {
		t.Errorf("Expected %s to be ''", nsQ.ToRequestParam())
	}
}
