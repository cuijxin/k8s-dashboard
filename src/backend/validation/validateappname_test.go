package validation

import (
	"testing"

	api "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestValidateName(t *testing.T) {
	spec := &AppNameValiditySpec{
		Namespace: "foo-namespace",
		Name:      "foo-name",
	}
	cases := []struct {
		spec     *AppNameValiditySpec
		objects  []runtime.Object
		expected bool
	}{
		{
			spec,
			nil,
			true,
		},
		{
			spec,
			[]runtime.Object{&api.ReplicationController{
				ObjectMeta: metaV1.ObjectMeta{
					Name:      "rc-1",
					Namespace: "ns-1",
				},
			}},
			true,
		},
		{
			spec,
			[]runtime.Object{&api.Service{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "rc-1", Namespace: "ns-1",
				},
			}},
			true,
		},
		{
			spec,
			[]runtime.Object{&api.ReplicationController{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "rc-1", Namespace: "ns-1",
				},
			}, &api.Service{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "rc-1", Namespace: "ns-1",
				},
			}},
			true,
		},
	}

	for _, c := range cases {
		testClient := fake.NewSimpleClientset(c.objects...)
		validity, _ := ValidateAppName(c.spec, testClient)
		if validity.Valid != c.expected {
			t.Errorf("Expected %#v validity to be %#v for objects %#v, but was %#v\n",
				c.spec, c.expected, c.objects, validity)
		}
	}
}
