package sidecar

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	metricapi "github.com/cuijxin/k8s-dashboard/src/backend/integration/metric/api"
)

// SidecarAllInOneDownloadConfig holds config information specifying whether given native Sidecar
// resource type supports list download.
var SidecarAllInOneDownloadConfig = map[api.ResourceKind]bool{
	api.ResourceKindPod:  true,
	api.ResourceKindNode: false,
}

// DataPointsFromMetricJSONFormat converts all the data points from format used by sidecar to our
// format.
func DataPointsFromMetricJSONFormat(raw []metricapi.MetricPoint) (dp metricapi.DataPoints) {
	for _, point := range raw {
		converted := metricapi.DataPoint{
			X: point.Timestamp.Unix(),
			Y: int64(point.Value),
		}

		if converted.Y < 0 {
			converted.Y = 0
		}

		dp = append(dp, converted)
	}
	return
}
