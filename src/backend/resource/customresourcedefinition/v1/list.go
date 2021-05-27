package v1

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/customresourcedefinition/types"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// GetCustomResourceDefinitionList returns all the custom resource definitions in the cluster.
func GetCustomResourceDefinitionList(
	client apiextensionsclientset.Interface,
	dsQuery *dataselect.DataSelectQuery) (*types.CustomResourceDefinitionList, error) {

	channel := common.GetCustomResourceDefinitionChannelV1(client, 1)
	crdList := <-channel.List
	err := <-channel.error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toCustomResoruceDefinitionList(crdList.Items, nonCriticalErrors, dsQuery), nil
}
