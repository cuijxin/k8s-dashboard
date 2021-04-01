package plugin

import (
	"net/http"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
)

// Config holds the information required by the frontend applicaiton to bootstrap.
type Config struct {
	Status         int32      `json:"status"`
	PluginMetadata []Metadata `json:"plugins"`
	Errors         []error    `json:"errors,omitempty"`
}

// Metadata holds the least possible plugin information for Config.
type Metadata struct {
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Dependenccies []string `json:"dependencies"`
}

func toPluginMetadata(vs []Plugin, f func(plugin Plugin) Metadata) []Metadata {
	vsm := make([]Metadata, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func statusCodeFromError(err error) int32 {
	if statusError, ok := err.(*apiErrors.StatusError); ok {
		return statusError.Status().Code
	}
	return http.StatusUnprocessableEntity
}
