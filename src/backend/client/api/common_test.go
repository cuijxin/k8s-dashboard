package api

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/authorization/v1"
)

func TestToSelfSubjectAccessReview(t *testing.T) {
	ns := "test-ns"
	name := "test-name"
	resourceName := "deployment"
	verb := "GET"
	expected := &v1.SelfSubjectAccessReview{
		Spec: v1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &v1.ResourceAttributes{
				Namespace: ns,
				Name:      name,
				Resource:  "deployments",
				Verb:      "get",
			},
		},
	}

	got := ToSelfSubjectAccessReview(ns, name, resourceName, verb)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected to get %+v but got %+v", expected, got)
	}
}
