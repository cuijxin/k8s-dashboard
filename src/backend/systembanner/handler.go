package systembanner

import (
	"net/http"

	"github.com/cuijxin/k8s-dashboard/src/backend/systembanner/api"
	restful "github.com/emicklei/go-restful/v3"
)

// SystemBannerHandler manages all endpoints related to system banner management.
type SystemBannerHandler struct {
	manager SystemBannerManager
}

// Install creates new endpoints for system banner management.
func (s *SystemBannerHandler) Install(ws *restful.WebService) {
	ws.Route(
		ws.GET("/systembanner").
			To(s.handleGet).
			Writes(api.SystemBanner{}))
}

func (s *SystemBannerHandler) handleGet(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, s.manager.Get())
}

// NewSystemBannerHandler creates SystemBannerHandler.
func NewSystemBannerHandler(manager SystemBannerManager) SystemBannerHandler {
	return SystemBannerHandler{manager: manager}
}
