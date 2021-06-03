package common

import (
	"reflect"
	"testing"

	api "k8s.io/api/core/v1"
)

func getReplicasPointer(replicas int32) *int32 {
	return &replicas
}

func TestGetPodInfo(t *testing.T) {
	cases := []struct {
		current  int32
		desired  *int32
		pods     []api.Pod
		expected PodInfo
	}{
		{
			5,
			getReplicasPointer(4),
			[]api.Pod{
				{
					Status: api.PodStatus{
						Phase: api.PodRunning,
					},
				},
			},
			PodInfo{
				Current:  5,
				Desired:  getReplicasPointer(4),
				Running:  1,
				Pending:  0,
				Failed:   0,
				Warnings: make([]Event, 0),
			},
		},
	}

	for _, c := range cases {
		actual := GetPodInfo(c.current, c.desired, c.pods)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("getPodInfo(%#v, %#v, %#v) == \n%#v\nexpected \n%#v\n", c.current, c.desired, c.pods, actual, c.expected)
		}
	}
}
