package integration

import (
	"net/http"

	"github.com/cuijxin/k8s-dashboard/src/backend/integration/api"
	restful "github.com/emicklei/go-restful/v3"
)

// IntegrationHandler messages all endpoints related to integrated applications, such as state.
type IntegrationHandler struct {
	manager IntegrationManager
}

// Install creates new endpoints for integrations. All information that any integration would want
// to expose by creating new endpoints should be kept here, i.e. helm integration might want to
// create endpoint to list available releases/charts.
//
// By default endpoint for checking state of the integrations is installed. It allows user
// to check state of integration by accessing `<DASHBOARD_URL>/api/v1/integration/{name}/state`.
func (self IntegrationHandler) Install(ws *restful.WebService) {
	ws.Route(
		ws.GET("/integration/{name}/state").
			To(self.handleGetState).
			Writes(api.IntegrationState{}))
}

func (self IntegrationHandler) handleGetState(request *restful.Request, response *restful.Response) {
	integrationName := request.PathParameter("name")
	state, err := self.manager.GetState(api.IntegrationID(integrationName))
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error()+"\n")
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, state)
}

// NewIntegrationHandler creates IntegrationHandler.
func NewIntegrationHandler(manager IntegrationManager) IntegrationHandler {
	return IntegrationHandler{manager: manager}
}
