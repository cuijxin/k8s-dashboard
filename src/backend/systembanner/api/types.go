package api

// SystemBannerManager is used for user system banner management.
type SystemBannerManager interface {
	// Get system banner.
	Get() *SystemBanner
}

// SystemBanner represents system banner.
type SystemBanner struct {
	Message  string               `json:"message"`
	Severity SystemBannerSeverity `json:"severity"`
}

// SystemBannerSeverity represents severity of system banner.
type SystemBannerSeverity string

const (
	// SystemBannerSeverityInfo is the lowest of allowed system banner severities.
	SystemBannerSeverityInfo SystemBannerSeverity = "INFO"

	// SystemBannerSeverityWarning is in the middle of allowed system banner severities.
	SystemBannerSeverityWarning SystemBannerSeverity = "WARNING"

	// SystemBannerSeverityError is the highest of allowed system banner severities.
	SystemBannerSeverityError SystemBannerSeverity = "ERROR"
)

// GetSeverity returns one of allowed severity values based on given parameter.
func GetSeverity(severity string) SystemBannerSeverity {
	switch severity {
	case string(SystemBannerSeverityWarning):
		return SystemBannerSeverityWarning
	case string(SystemBannerSeverityError):
		return SystemBannerSeverityError
	default:
		return SystemBannerSeverityInfo
	}
}
