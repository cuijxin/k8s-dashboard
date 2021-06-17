package handler

import (
	"net/http"

	"github.com/cuijxin/k8s-dashboard/src/backend/auth"
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/handler/parser"
	"github.com/cuijxin/k8s-dashboard/src/backend/integration"
	"github.com/cuijxin/k8s-dashboard/src/backend/plugin"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/clusterrole"
	"github.com/cuijxin/k8s-dashboard/src/backend/settings"
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

	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	integrationHandler := integration.NewIntegrationHandler(iManager)
	integrationHandler.Install(apiV1Ws)

	plguinHandler := plugin.NewPluginHandler(cManager)
	plguinHandler.Install(apiV1Ws)

	authHandler := auth.NewAuthHandler(authManager)
	authHandler.Install(apiV1Ws)

	settingsHandler := settings.NewSettingsHandler(sManager, cManager)
	settingsHandler.Install(apiV1Ws)

	systemBannerHandler := systembanner.NewSystemBannerHandler(sbManager)
	systemBannerHandler.Install(apiV1Ws)

	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrole").
			To(apiHandler.handleGetClusterRoleList).
			Writes(clusterrole.ClusterRoleList{}))

	return wsContainer, nil
}

func (apiHandler *APIHandler) handleGetClusterRoleList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrole.GetClusterRoleList(k8sClient, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
