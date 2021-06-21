package role

import (
	"log"

	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

// GetRoleList returns a list of all Roles in the cluster.
func GetRoleList(client kubernetes.Interface, nsQuery *common.NamespaceQuery, dsQuery *dataselect.DataSelectQuery) (*RoleList, error) {
	log.Print("Getting list of all roles in the cluster")
	channels := &common.ResourceChannels{
		RoleList: common.GetRoleListChannel(client, nsQuery, 1),
	}

	return GetRoleListFromChannels(channels, dsQuery)
}

// GetRoleListFromChannels returns a list of all Roles in the cluster
// reading required resource list once from the channels.
func GetRoleListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (*RoleList, error) {
	roles := <-channels.RoleList.List
	err := <-channels.RoleList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}
	roleList := toRoleList(roles.Items, nonCriticalErrors, dsQuery)
	return roleList, nil
}

func toRole(role rbac.Role) Role {
	return Role{
		ObjectMeta: api.NewObjectMeta(role.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindRole),
	}
}

func toRoleList(roles []rbac.Role, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *RoleList {
	result := &RoleList{
		ListMeta: api.ListMeta{TotalItems: len(roles)},
		Errors:   nonCriticalErrors,
	}

	items := make([]Role, 0)
	for _, item := range roles {
		items = append(items, toRole(item))
	}

	roleCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(items), dsQuery)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	result.Items = fromCells(roleCells)
	return result
}
