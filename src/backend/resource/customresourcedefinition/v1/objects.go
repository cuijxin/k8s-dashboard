package v1

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/customresourcedefinition/types"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
)

// GetCustomResourceObjectList gets objects for a CR.
func GetCustomResourceObjectList(
	client apiextensionsclientset.Interface,
	config *rest.Config,
	namespace *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery,
	crdName string) (*types.CustomResourceObjectList, error) {

}
