package sidecar

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/cuijxin/k8s-dashboard/src/backend/args"
)

// SidecarRESTClient is used to make raw requests to sidecar.
type SidecarRESTClient interface {
	// Creates a new GET HTTP request to sidecar, specified by the path param, to the V1 API
	// endpoint. The path param is without the API prefix, e.g.,
	// /model/namespaces/default/pod-list/foo/metrics/memory-usage
	Get(path string) RequestInterface
	HealthCheck() error
}

// RequestInterface is an interface that allows to make operations on pure request object.
// Separation is done to allow testing.
type RequestInterface interface {
	DoRaw(context.Context) ([]byte, error)
	AbsPath(segments ...string) *rest.Request
}

// InClusterSidecarClient is an in-cluster implementation of a Sidecar client. Talks with Sidecar
// through service proxy.
type inClusterSidecarClient struct {
	client rest.Interface
}

// Get creates request to given path.
func (c inClusterSidecarClient) Get(path string) RequestInterface {
	return c.client.Get().
		Namespace(args.Holder.GetNamespace()).
		Resource("services").
		Name("dashboard-metrics-scraper").
		SubResource("proxy").
		Suffix(path)
}

// HealthCheck does a health check of the application.
// Returns nil if connection to application can be established, error object otherwise.
func (self inClusterSidecarClient) HealthCheck() error {
	_, err := self.client.Get().
		Namespace(args.Holder.GetNamespace()).
		Resource("services").
		Name("dashboard-metrics-scraper").
		SubResource("proxy").
		Suffix("/healthz").
		DoRaw(context.TODO())
	return err
}

// RemoteSidecarClient is an implementation of a remote Sidecar client. Talks with Sidecar
// through raw RESTClient.
type remoteSidecarClient struct {
	client rest.Interface
}

// Get creates request to given path.
func (c remoteSidecarClient) Get(path string) RequestInterface {
	return c.client.Get().Suffix(path)
}

// HealthCheck does a health check of the application.
// Returns nil if connection to application can be established, error object otherwise.
func (self remoteSidecarClient) HealthCheck() error {
	_, err := self.Get("healthz").AbsPath("/").DoRaw(context.TODO())
	return err
}
