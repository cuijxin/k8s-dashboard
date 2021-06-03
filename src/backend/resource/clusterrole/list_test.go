package clusterrole

import (
	"reflect"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	rbac "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestToRbacRoleLists(t *testing.T) {
	cases := []struct {
		clusterRoles []rbac.ClusterRole
		expected     *ClusterRoleList
	}{
		{nil, &ClusterRoleList{Items: []ClusterRole{}}},
		{
			[]rbac.ClusterRole{
				{
					ObjectMeta: metaV1.ObjectMeta{Name: "cluster-role"},
					Rules: []rbac.PolicyRule{{
						Verbs:     []string{"post", "put"},
						Resources: []string{"pods", "deployments"},
					}},
				},
			},
			&ClusterRoleList{
				ListMeta: api.ListMeta{TotalItems: 1},
				Items: []ClusterRole{{
					ObjectMeta: api.ObjectMeta{Name: "cluster-role", Namespace: ""},
					TypeMeta:   api.TypeMeta{Kind: api.ResourceKindClusterRole},
				}},
			},
		},
	}
	for _, c := range cases {
		actual := toClusterRoleLists(c.clusterRoles, nil, dataselect.NoDataSelect)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("toRbacRoleLists(%#v) == \n%#v\nexpected \n%#v\n",
				c.clusterRoles, actual, c.expected)
		}
	}
}
