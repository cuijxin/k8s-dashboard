package v1beta1

import (
	"reflect"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/customresourcedefinition/types"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetCustomResourceDefinition(t *testing.T) {
	cases := []struct {
		expectedActions []string
		crdList         *apiextensionsv1beta1.CustomResourceDefinitionList
		expected        *types.CustomResourceDefinitionList
	}{
		{
			[]string{"list"},
			&apiextensionsv1beta1.CustomResourceDefinitionList{
				Items: []apiextensionsv1beta1.CustomResourceDefinition{
					{
						ObjectMeta: metaV1.ObjectMeta{Name: "foos.samplecontroller.k8s.io"},
						Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
							Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
								Kind:   "Foo",
								Plural: "foos",
							},
							Versions: []apiextensionsv1beta1.CustomResourceDefinitionVersion{
								{
									Name: "v1alpha1",
								},
							},
						},
					},
				},
			},
			&types.CustomResourceDefinitionList{
				ListMeta: api.ListMeta{TotalItems: 1},
				Items: []types.CustomResourceDefinition{
					{
						ObjectMeta:  api.ObjectMeta{Name: "foos.samplecontroller.k8s.io"},
						TypeMeta:    api.TypeMeta{Kind: api.ResourceKindCustomResourceDefinition},
						Version:     "v1alpha1",
						Established: apiextensions.ConditionUnknown,
					},
				},
				Errors: []error{},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.crdList)

		actual, _ := GetCustomResourceDefinitionList(fakeClient, dataselect.DefaultDataSelect)

		actions := fakeClient.Actions()
		if len(actions) != len(c.expectedActions) {
			t.Errorf("Unexpected actions: %v, expected %d actions got %d", actions,
				len(c.expectedActions), len(actions))
			continue
		}

		for i, verb := range c.expectedActions {
			if actions[i].GetVerb() != verb {
				t.Errorf("Unexpected action: %+v, expected %s",
					actions[i], verb)
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetCustomResourceDefinitionList(client, nil) == \ngot: %#v, \nexpected %#v",
				actual, c.expected)
		}
	}
}
