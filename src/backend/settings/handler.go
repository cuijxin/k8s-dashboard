package settings

import (
	clientapi "github.com/cuijxin/k8s-dashboard/src/backend/client/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/settings/api"
)

// SettingsHandler manages all endpoints related to settings management.
type Settingshandler struct {
	manager       api.SettingsManager
	clientManager clientapi.ClientManager
}
