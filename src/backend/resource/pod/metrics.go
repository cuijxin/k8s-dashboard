package pod

import (
	"log"

	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	metricapi "github.com/cuijxin/k8s-dashboard/src/backend/integration/metric/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// MetricsByPod is a metrics map by pod name.
type MetricsByPod struct {
	// Metrics by namespace and name of a pod.
	MetricsMap map[types.UID]PodMetrics `json:"metricsMap"`
}

// PodMetrics is a structure representing pods metrics, contains information about
// CPU and memory usage.
type PodMetrics struct {
	// Most recent measure of CPU usage on all cores in nanoseconds.
	CPUUsage *uint64 `json:"cpuUsage"`
	// Pod memory usage in bytes.
	MemoryUsage *uint64 `json:"memoryUsage"`
	// Timestamped samples of CPUUsage over some short period of history
	CPUUsageHistory []metricapi.MetricPoint `json:"cpuUsageHistory"`
	// Timestamped samples of pod memory usage some short period of history
	MemoryUsageHistory []metricapi.MetricPoint `json:"memoryUsageHistory"`
}

func getMetricsPerPod(pods []v1.Pod, metricClient metricapi.MetricClient, dsQuery *dataselect.DataSelectQuery) (*MetricsByPod, error) {
	log.Println("Getting pod metrics")

	result := &MetricsByPod{MetricsMap: make(map[types.UID]PodMetrics)}

	metricPromises := dataselect.PodListMetrics(toCells(pods), dsQuery, metricClient)
	metrics, err := metricPromises.GetMetrics()
	if err != nil {
		return result, err
	}

	for _, m := range metrics {
		uid, err := getPodUIDFromMetric(m)
		if err != nil {
			log.Printf("Skipping metric because of error: %s", err.Error())
		}

		podMetrics := PodMetrics{}
		if p, exists := result.MetricsMap[uid]; exists {
			podMetrics = p
		}

		if m.MetricName == metricapi.CpuUsage && len(m.MetricPoints) > 0 {
			podMetrics.CPUUsage = &m.MetricPoints[len(m.MetricPoints)-1].Value
			podMetrics.CPUUsageHistory = m.MetricPoints
		}

		if m.MetricName == metricapi.MemoryUsage && len(m.MetricPoints) > 0 {
			podMetrics.MemoryUsage = &m.MetricPoints[len(m.MetricPoints)-1].Value
			podMetrics.MemoryUsageHistory = m.MetricPoints
		}

		result.MetricsMap[uid] = podMetrics
	}

	return result, nil
}

func getPodUIDFromMetric(metric metricapi.Metric) (types.UID, error) {
	// Check is metric label contains required resource UID
	uidList, exists := metric.Label[api.ResourceKindPod]
	if !exists {
		return "", errors.NewInvalid("Metric label not set.")
	}

	// Check if metric maps to single resource. Multiple uids means that data was aggregated
	// from multiple resources. We should have metrics per resource here.
	if len(uidList) != 1 {
		return "", errors.NewInvalid("Found multiple UIDs. Metric should contain data for single resource only.")
	}

	return uidList[0], nil
}
