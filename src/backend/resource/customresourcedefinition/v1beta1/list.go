package v1beta1

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/customresourcedefinition/types"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// GetCustomResourceDefinitionList returns all the custom resource definitions in the cluster.
func GetCustomResourceDefinitionList(client apiextensionsclientset.Interface,
	dsQuery *dataselect.DataSelectQuery) (*types.CustomResourceDefinitionList, error) {
	channel := common.GetCustomResourceDefinitionChannelV1beta1(client, 1)
	crdList := <-channel.List
	err := <-channel.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toCustomResourceDefinitionList(crdList.Items, nonCriticalErrors, dsQuery), nil
}

func toCustomResourceDefinitionList(crds []apiextensionsv1beta1.CustomResourceDefinition,
	nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *types.CustomResourceDefinitionList {
	crdList := &types.CustomResourceDefinitionList{
		Items:    make([]types.CustomResourceDefinition, 0),
		ListMeta: api.ListMeta{TotalItems: len(crds)},
		Errors:   nonCriticalErrors,
	}

	crdCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(crds), dsQuery)
	crds = fromCells(crdCells)
	crdList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, crd := range crds {
		crdList.Items = append(crdList.Items, toCustomResourceDefinition(&crd))
	}

	return crdList
}

func toCustomResourceDefinition(crd *apiextensionsv1beta1.CustomResourceDefinition) types.CustomResourceDefinition {
	return types.CustomResourceDefinition{
		ObjectMeta:  api.NewObjectMeta(crd.ObjectMeta),
		TypeMeta:    api.NewTypeMeta(api.ResourceKindCustomResourceDefinition),
		Version:     crd.Spec.Versions[0].Name,
		Group:       crd.Spec.Group,
		Scope:       toCustomResourceDefinitionScope(crd.Spec.Scope),
		Names:       toCustomResourceDefinitionAcceptedNames(crd.Status.AcceptedNames),
		Established: getCRDConditionStatus(crd, apiextensionsv1beta1.Established),
	}
}

func toCustomResourceDefinitionScope(scope apiextensionsv1beta1.ResourceScope) apiextensions.ResourceScope {
	return apiextensions.ResourceScope(scope)
}

func toCustomResourceDefinitionAcceptedNames(names apiextensionsv1beta1.CustomResourceDefinitionNames) types.CustomResourceDefinitionNames {
	return types.CustomResourceDefinitionNames{
		Plural:     names.Plural,
		Singular:   names.Singular,
		ShortNames: names.ShortNames,
		Kind:       names.Kind,
		ListKind:   names.ListKind,
		Categories: names.Categories,
	}
}

func getCRDConditionStatus(node *apiextensionsv1beta1.CustomResourceDefinition, conditionType apiextensionsv1beta1.CustomResourceDefinitionConditionType) apiextensions.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return apiextensions.ConditionStatus(condition.Status)
		}
	}
	return apiextensions.ConditionUnknown
}
