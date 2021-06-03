package common

import (
	v1 "k8s.io/api/core/v1"

	"github.com/cuijxin/k8s-dashboard/src/backend/api"
)

// FilterNamespacedServicesBySelector returns services targeted by given resource selector in
// given namespace.
func FilterNamespacedServicesBySelector(services []v1.Service, namespace string, resourceSelector map[string]string) []v1.Service {
	var matchingServices []v1.Service
	for _, service := range services {
		if service.ObjectMeta.Namespace == namespace && api.IsSelectorMatching(service.Spec.Selector, resourceSelector) {
			matchingServices = append(matchingServices, service)
		}
	}

	return matchingServices
}
