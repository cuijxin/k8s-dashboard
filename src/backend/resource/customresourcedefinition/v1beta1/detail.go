package v1beta1

import (
	"context"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/customresourcedefinition/types"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// GetCustomResourceDefinitionDetail returns detailed information about custom resource definition.
func GetCustomResourceDefinitionDetail(client apiextensionsclientset.Interface,
	config *rest.Config, name string) (*types.CustomResourceDefinitionDetail, error) {
	customResourceDefinition, err := client.ApiextensionsV1beta1().CustomResourceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	objects, err := GetCustomResourceObjectList(client, config, &common.NamespaceQuery{}, dataselect.DefaultDataSelect, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	return toCustomResourceDefinitionDetail(customResourceDefinition, *objects, nonCriticalErrors), nil
}

func toCustomResourceDefinitionDetail(crd *apiextensions.CustomResourceDefinition,
	objects types.CustomResourceObjectList,
	nonCriticalErrors []error) *types.CustomResourceDefinitionDetail {
	subresources := []string{}
	crdSubresources := crd.Spec.Versions[0].Subresources
	if crdSubresources != nil {
		if crdSubresources != nil {
			subresources = append(subresources, "Scale")
		}
		if crdSubresources.Status != nil {
			subresources = append(subresources, "Status")
		}
	}

	return &types.CustomResourceDefinitionDetail{
		CustomResourceDefinition: toCustomResourceDefinition(crd),
		Versions:                 getCRDVersions(crd),
		Conditions:               getCRDConditions(crd),
		Objects:                  objects,
		Subresources:             subresources,
		Errors:                   nonCriticalErrors,
	}
}

func getCRDVersions(crd *apiextensions.CustomResourceDefinition) []types.CustomResourceDefinitionVersion {
	crdVersions := make([]types.CustomResourceDefinitionVersion, 0, len(crd.Spec.Versions))
	if len(crd.Spec.Versions) > 0 {
		for _, version := range crd.Spec.Versions {
			crdVersions = append(crdVersions, types.CustomResourceDefinitionVersion{
				Name:    version.Name,
				Served:  version.Served,
				Storage: version.Storage,
			})
		}
	}

	return crdVersions
}
