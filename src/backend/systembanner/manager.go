package systembanner

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/systembanner/api"
)

// SystemBannerManager is a structure containing all system banner manager members.
type SystemBannerManager struct {
	systemBanner api.SystemBanner
}

// NewSystemBannerManager creates new settings manager.
func NewSystemBannerManager(message, severity string) SystemBannerManager {
	return SystemBannerManager{
		systemBanner: api.SystemBanner{
			Message:  message,
			Severity: api.GetSeverity(severity),
		},
	}
}

// Get implements SystemBannerManager interface. Check if for more information.
func (sbm *SystemBannerManager) Get() api.SystemBanner {
	return sbm.systemBanner
}
