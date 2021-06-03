package common

import (
	"reflect"
	"testing"

	batch "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type metaObj struct {
	metaV1.ObjectMeta
	metaV1.TypeMeta
}

func TestFilterPodsByControllerRef(t *testing.T) {
	controller := true
	okOwnerRef := []metaV1.OwnerReference{{
		Kind:       "ReplicationController",
		Name:       "my-name-1",
		UID:        "uid-1",
		Controller: &controller,
	}}
	nokOwnerRef := []metaV1.OwnerReference{{
		Kind:       "ReplicationController",
		Name:       "my-name-1",
		UID:        "",
		Controller: &controller,
	}}
	cases := []struct {
		obj      *metaObj
		pods     []api.Pod
		expected []api.Pod
	}{
		{
			&metaObj{
				ObjectMeta: metaV1.ObjectMeta{
					UID:  "uid-1",
					Name: "my-name-1",
				},
			},
			[]api.Pod{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:            "first-pod-ok",
						Namespace:       "default",
						OwnerReferences: okOwnerRef,
					},
				},
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:            "second-pod-ok",
						Namespace:       "default",
						OwnerReferences: okOwnerRef,
					},
				},
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:            "third-pod-nok",
						Namespace:       "default",
						OwnerReferences: nokOwnerRef,
					},
				},
			},
			[]api.Pod{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:            "first-pod-ok",
						Namespace:       "default",
						OwnerReferences: okOwnerRef,
					},
				},
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:            "second-pod-ok",
						Namespace:       "default",
						OwnerReferences: okOwnerRef,
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual := FilterPodsByControllerRef(c.obj, c.pods)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("FilterPodsByControllerRef(%+v, %+v) == %+v, expected %+v", c.pods, c.obj, actual, c.expected)
		}
	}
}

func TestGetContainerImages(t *testing.T) {
	cases := []struct {
		podTemplate *api.PodSpec
		expected    []string
	}{
		{&api.PodSpec{}, nil},
		{
			&api.PodSpec{
				Containers: []api.Container{{Image: "container:v1"}, {Image: "container:v2"}},
			},
			[]string{"container:v1", "container:v2"},
		},
	}

	for _, c := range cases {
		actual := GetContainerImages(c.podTemplate)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetContainerImages(%+v) == %+v, expected %+v", c.podTemplate, actual, c.expected)
		}
	}
}

func TestGetInitContainerImages(t *testing.T) {
	cases := []struct {
		podTemplate *api.PodSpec
		expected    []string
	}{
		{&api.PodSpec{}, nil},
		{
			&api.PodSpec{
				InitContainers: []api.Container{{Image: "initContainer:v3"}, {Image: "initContainer:v4"}},
			},
			[]string{"initContainer:v3", "initContainer:v4"},
		},
	}

	for _, c := range cases {
		actual := GetInitContainerImages(c.podTemplate)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetInitContainerImages(%+v) == %+v, expected %+v", c.podTemplate, actual, c.expected)
		}
	}
}

func TestGetContainerNames(t *testing.T) {
	cases := []struct {
		podTemplate *api.PodSpec
		expected    []string
	}{
		{&api.PodSpec{}, nil},
		{
			&api.PodSpec{
				Containers: []api.Container{{Name: "container"}, {Name: "container"}},
			},
			[]string{"container", "container"},
		},
	}

	for _, c := range cases {
		actual := GetContainerNames(c.podTemplate)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetContainerNames(%+v) == %+v, expected %+v",
				c.podTemplate, actual, c.expected)
		}
	}
}

func TestGetInitContainerNames(t *testing.T) {
	cases := []struct {
		podTemplate *api.PodSpec
		expected    []string
	}{
		{&api.PodSpec{}, nil},
		{
			&api.PodSpec{
				InitContainers: []api.Container{{Name: "initContainer"}, {Name: "initContainer"}},
			},
			[]string{"initContainer", "initContainer"},
		},
	}

	for _, c := range cases {
		actual := GetInitContainerNames(c.podTemplate)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetContainerNames(%+v) == %+v, expected %+v",
				c.podTemplate, actual, c.expected)
		}
	}
}

func TestGetNonduplicateContaienrImages(t *testing.T) {
	expected := []string{"Container1", "Container2", "Container3"}
	pods := make([]api.Pod, 2, 2)

	pods[0] = api.Pod{
		Spec: api.PodSpec{
			Containers: []api.Container{{Image: "Container1"}, {Image: "Container2"}},
		},
	}

	pods[1] = api.Pod{
		Spec: api.PodSpec{
			Containers: []api.Container{{Image: "Container2"}, {Image: "Container3"}},
		},
	}

	actual := GetNonduplicateContainerImages(pods)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetNonduplicateContainerImages() == %+v, expected %+v", actual, expected)
	}
}

func TestGetNonduplicateInitContainerImages(t *testing.T) {
	expected := []string{"initContainer1", "initContainer2", "initContainer3"}
	pods := make([]api.Pod, 2, 2)

	pods[0] = api.Pod{
		Spec: api.PodSpec{
			InitContainers: []api.Container{{Image: "initContainer1"}, {Image: "initContaienr2"}},
		},
	}

	pods[1] = api.Pod{
		Spec: api.PodSpec{
			InitContainers: []api.Container{{Image: "initContainer2"}, {Image: "initContainer3"}},
		},
	}
	actual := GetNonduplicateInitContainerImages(pods)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetNonduplicateInitContainerImages() == %+v, expected %+v", actual, expected)
	}
}

func TestGetNonduplicateInitContainerNames(t *testing.T) {
	expected := []string{"initContainer1", "initContainer2", "initContainer3"}
	pods := make([]api.Pod, 2, 2)

	pods[0] = api.Pod{
		Spec: api.PodSpec{
			InitContainers: []api.Container{{Name: "initContainer1"}, {Name: "initContainer2"}},
		},
	}

	pods[1] = api.Pod{
		Spec: api.PodSpec{
			InitContainers: []api.Container{{Name: "initContainer2"}, {Name: "initContainer3"}},
		},
	}
	actual := GetNonduplicateInitContainerNames(pods)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetNonduplicateInitContainerNames() == %+v, expected %+v", actual, expected)
	}
}

func TestFilterPodsForJob(t *testing.T) {
	cases := []struct {
		job      batch.Job
		pods     []api.Pod
		expected []api.Pod
	}{
		{
			batch.Job{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "default",
					Name:      "job-1",
					UID:       "job-uid",
				},
				Spec: batch.JobSpec{
					Selector: &metaV1.LabelSelector{
						MatchLabels: map[string]string{"controller-uid": "job-uid"},
					},
				},
			},
			[]api.Pod{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "pod-1",
						Namespace: "default",
						Labels:    map[string]string{"controller-uid": "job-uid"},
					},
				},
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "pod-2",
						Namespace: "default",
					},
				},
			},
			[]api.Pod{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "pod-1",
						Namespace: "default",
						Labels:    map[string]string{"controller-uid": "job-uid"},
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual := FilterPodsForJob(c.job, c.pods)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("FilterPodsForJob(%+v, %+v) == %+v, expected %+v", c.job, c.pods, actual, c.expected)
		}
	}
}
