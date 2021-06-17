package validation

import (
	"reflect"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"

	test "k8s.io/client-go/testing"
)

func areErrorsEqual(err1, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() == err2.Error()) || (err1 == nil && err2 == nil)
}

type fakeServerGroupsMethod func() (*metav1.APIGroupList, error)

type FakeClient struct {
	fake.Clientset
	fakeServerGroupsMethod
}

func (f *FakeClient) Discovery() discovery.DiscoveryInterface {
	return &FakeDiscovery{
		FakeDiscovery:          fakediscovery.FakeDiscovery{Fake: &f.Fake},
		fakeServerGroupsMethod: f.fakeServerGroupsMethod,
	}
}

type FakeDiscovery struct {
	fakediscovery.FakeDiscovery
	fakeServerGroupsMethod
}

func (f *FakeDiscovery) ServerGroups() (*metav1.APIGroupList, error) {
	return f.fakeServerGroupsMethod()
}

func TestValidateRbacStatus(t *testing.T) {
	cases := []struct {
		info        string
		mockMethod  fakeServerGroupsMethod
		expected    *RbacStatus
		expectedErr error
	}{
		{
			"should throw an error when can't get api versions from server",
			func() (*metav1.APIGroupList, error) {
				return nil, errors.NewInvalid("test-error")
			},
			nil,
			errors.NewInvalid("Couldn't get available api versions from server: test-error"),
		},
		{
			"should disable rbacs when supported api version not enabled on the server",
			func() (*metav1.APIGroupList, error) {
				return &metav1.APIGroupList{Groups: []metav1.APIGroup{
					{Name: "rbac", Versions: []metav1.GroupVersionForDiscovery{
						{
							GroupVersion: "authorization.k8s.io/v1alpha1",
							Version:      "v1alpha1",
						},
					}},
				}}, nil
			},
			&RbacStatus{false},
			nil,
		},
		{
			"should enable rbac when supported api version is enabled on the server",
			func() (*metav1.APIGroupList, error) {
				return &metav1.APIGroupList{Groups: []metav1.APIGroup{
					{Name: "rbac", Versions: []metav1.GroupVersionForDiscovery{
						{
							GroupVersion: "authorization.k8s.io/v1beta1",
							Version:      "v1beta1",
						},
					}},
				}}, nil
			},
			&RbacStatus{true},
			nil,
		},
	}

	for _, c := range cases {
		client := &FakeClient{
			Clientset:              fake.Clientset{Fake: test.Fake{}},
			fakeServerGroupsMethod: c.mockMethod,
		}

		status, err := ValidateRbacStatus(client)
		if !areErrorsEqual(err, c.expectedErr) {
			t.Fatalf("Test case: %s. Expected error to be: %v, but got %v.", c.info,
				c.expectedErr, err)
		}

		if !reflect.DeepEqual(status, c.expected) {
			t.Fatalf("Test case: %s. Expected status to be: %v, but got %v.", c.info,
				c.expected, status)
		}
	}

}
