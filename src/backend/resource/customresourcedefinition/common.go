package customresourcedefinition

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
)

var (
	groupName = apiextensionsv1.GroupName
	v1        = apiextensionsv1.SchemeGroupVersion.Version
	v1beta1   = apiextensionsv1beta.SchemeGroupVersion.Version
)

func GetExtensionsAPIVersion(client clientset.Interface) (string, error) {
	list, err := client.Discovery().ServerGroups()
	if err != nil {
		return "", err
	}

	for _, group := range list.Groups {
		if group.Name == groupName {
			return group.PreferredVersion.Version, nil
		}
	}

	return "", errors.NewNotFound("supported version for extensions api not found")
}
