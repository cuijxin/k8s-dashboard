package handler

import (
	"net/http"

	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration"
	settingsApi "github.com/cuijxin/k8s-dashboard/src/backend/settings/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/systembanner"
	"github.com/emicklei/go-restful/v3"
)

const (
	// RequestLogString is a template for request log message.
	RequestLogString = "[%s] Incoming %s %s %s request from %s: %s"

	// ResponseLogString is a template for response log message.
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

// APIHandler is a representation of API handler. Structure contains clientapi,
// Heapster clientapi and clientapi configuration.
type APIHandler struct {
	iManager integration.IntegrationManager
	cManager clientapi.ClientManager
	sManager settingsApi.SettingsManager
}

// TerminalResponse is sent by handleExecShell. The Id is a random session id
// that binds the original REST request and the SockJS connection.
// Any clientapi in possession of this Id can hijack the terminal session.
type TerminalResponse struct {
	ID string `json:"id"`
}

// CreateHTTPAPIHandler creates a new HTTP handler that handles all requests to
// the API of the backend.
func CreateHTTPAPIHandler(iManager integration.IntegrationManager,
	cManager clientapi.ClientManager,
	authManager authApi.AuthManager,
	sManager settingsApi.SettingsManager,
	sbManager systembanner.SystemBannerManager) (http.Handler, error) {
	apiHandler := APIHandler{iManager: iManager, cManager: cManager, sManager: sManager}
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)

	InstallFilters(apiV1Ws, cManager)
}
