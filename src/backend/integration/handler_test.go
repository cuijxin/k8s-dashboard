package integration

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestIntegrationHandler_Install(t *testing.T) {
	iHandler := NewIntegrationHandler(nil)
	ws := new(restful.WebService)
	iHandler.Install(ws)

	if len(ws.Routes()) == 0 {
		t.Error("Failed to install routes.")
	}
}
