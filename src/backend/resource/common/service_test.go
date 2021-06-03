package common

import (
	"reflect"
	"testing"

	api "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFilterNamespacedServicesBySelector(t *testing.T) {
	firstLabelSelectorMap := make(map[string]string)
	firstLabelSelectorMap["name"] = "app-name-first"
	secondLabelSelectorMap := make(map[string]string)
	secondLabelSelectorMap["name"] = "app-name-second"

	cases := []struct {
		selector  map[string]string
		namespace string
		services  []api.Service
		expected  []api.Service
	}{
		{
			firstLabelSelectorMap, "test-ns-1",
			[]api.Service{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "first-service-ok",
						Namespace: "test-ns-1",
					},
					Spec: api.ServiceSpec{
						Selector: firstLabelSelectorMap,
					},
				},
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "second-service-wrong",
						Namespace: "test-ns-2",
					},
					Spec: api.ServiceSpec{
						Selector: firstLabelSelectorMap,
					},
				},
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "third-service-wrong",
						Namespace: "test-ns-1",
					},
					Spec: api.ServiceSpec{
						Selector: secondLabelSelectorMap,
					},
				},
			},
			[]api.Service{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "first-service-ok",
						Namespace: "test-ns-1",
					},
					Spec: api.ServiceSpec{
						Selector: firstLabelSelectorMap,
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual := FilterNamespacedServicesBySelector(c.services, c.namespace, c.selector)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("FilterNamespacedServicesBySelector(%+v, %+v) == \n%+v, expected \n%+v",
				c.services, c.selector, actual, c.expected)
		}
	}
}
