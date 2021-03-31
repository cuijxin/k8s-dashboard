package v1

import (
	"context"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/customresourcedefinition/types"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
)

// GetCustomResourceDefinitionDetail returns detailed information about a custom
// resource definition.
func GetCustomResourceDefinitionDetail(client apiextensionsclientset.Interface, config *rest.Config, name string) (*types.CustomResourceDefinitionDetail, error) {
	customResourceDefinition, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
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
