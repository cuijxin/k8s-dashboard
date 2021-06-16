package plugin

import (
	"net/http"

	"github.com/cuijxin/k8s-dashboard/src/backend/handler/parser"
	"github.com/emicklei/go-restful/v3"
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
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies"`
}

func toPluginMetadata(vs []Plugin, f func(plugin Plugin) Metadata) []Metadata {
	vsm := make([]Metadata, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func (h *Handler) handleConfig(request *restful.Request, response *restful.Response) {
	pluginClient, err := h.cManager.PluginClient(request)
	cfg := Config{Status: http.StatusOK, PluginMetadata: []Metadata{}, Errors: []error{}}
	if err != nil {
		cfg.Status = statusCodeFromError(err)
		cfg.Errors = append(cfg.Errors, err)
		response.WriteHeaderAndEntity(http.StatusOK, cfg)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := GetPluginList(pluginClient, "", dataSelect)
	if err != nil {
		cfg.Status = statusCodeFromError(err)
		cfg.Errors = append(cfg.Errors, err)
		response.WriteHeaderAndEntity(http.StatusOK, cfg)
		return
	}

	if result != nil && len(result.Errors) > 0 {
		cfg.Status = statusCodeFromError(result.Errors[0])
		cfg.Errors = append(cfg.Errors, result.Errors...)
		response.WriteHeaderAndEntity(http.StatusOK, cfg)
		return
	}

	cfg.PluginMetadata = toPluginMetadata(result.Items, func(plugin Plugin) Metadata {
		return Metadata{
			Name:         plugin.Name,
			Path:         plugin.Path,
			Dependencies: plugin.Dependencies,
		}
	})
	cfg.Errors = result.Errors
	response.WriteHeaderAndEntity(http.StatusOK, cfg)
}

func statusCodeFromError(err error) int32 {
	if statusError, ok := err.(*apiErrors.StatusError); ok {
		return statusError.Status().Code
	}
	return http.StatusUnprocessableEntity
}
