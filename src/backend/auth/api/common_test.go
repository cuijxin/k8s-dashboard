package api

import (
	"reflect"
	"testing"
)

func TestToAuthenticationModes(t *testing.T) {
	cases := []struct {
		modes    []string
		expected AuthenticationModes
	}{
		{[]string{}, AuthenticationModes{}},
		{[]string{"token"}, AuthenticationModes{Token: true}},
		{[]string{"token", "basic", "test"}, AuthenticationModes{Token: true, Basic: true}},
	}

	for _, c := range cases {
		got := ToAuthenticationModes(c.modes)
		if !reflect.DeepEqual(got, c.expected) {
			t.Fatalf("ToAuthenticationModes(): expected %v, but got %v", c.expected, got)
		}
	}
}

func TestShouldRejectRequest(t *testing.T) {
  cases := []struct {
    url string
    expected bool
  }{
    {"#!/namespace?namespace=test", false},
    {"#!/namespace/test", false},
		{"#!/namespace?namespace=kube-system", false},
		{"#!/namespace/kube-system", false},
		{"#!/secret/test/test-secret?namespace=test", false},
		{"#!/secret/kube-system/test-secret", false},
		{"#!/secret/kube-system/kubernetes-dashboard-key-holder", true},
		{"#!/secret/test/kubernetes-dashboard-certs", true},
		{"#!/secret/kube-system/kubernetes-dashboard-certs", true},
	}

	for _, c := range cases {
		got := ShouldRejectRequest(c.url)
		if !reflect.DeepEqual(got, c.expected) {
			t.Fatalf("ShouldRejectRequest(): url %s expected %v, but got %v", c.url, c.expected, got)
		}
	}
}
