package event

import (
	"context"

	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

// EmptyEventList is a empty list of events.
var EmptyEventList = &common.EventList{
	Events: make([]common.Event, 0),
	ListMeta: api.ListMeta{
		TotalItems: 0,
	},
}

// GetEvents gets events associated to resoruce with given name.
func GetEvents(client kubernetes.Interface, namespace, resourceName string) ([]v1.Event, error) {
	fieldSelector, err := fields.ParseSelector("involvedObject.name" + "=" + resourceName)

	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		EventList: common.GetEventListChannelWithOptions(
			client,
			common.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labels.Everything().String(),
				FieldSelector: fieldSelector.String(),
			},
			1),
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	return FillEventsType(eventList.Items), nil
}

// GetPodsEvents gets events targeting given list of pods.
func GetPodsEvents(client kubernetes.Interface, namespace string, pods []v1.Pod) ([]v1.Event, error) {
	nsQuery := common.NewSameNamespaceQuery(namespace)
	if namespace == v1.NamespaceAll {
		nsQuery = common.NewNamespaceQuery([]string{})
	}

	channels := &common.ResourceChannels{
		EventList: common.GetEventListChannel(client, nsQuery, 1),
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	events := filterEventsByPodsUID(eventList.Items, pods)

	return events, nil
}

// GetPodEvents gets pods events associated to pod name and namespace
func GetPodEvents(client kubernetes.Interface, namespace, podName string) ([]v1.Event, error) {
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannel(client,
			common.NewSameNamespaceQuery(namespace),
			1),
		EventList: common.GetEventListChannel(client, common.NewSameNamespaceQuery(namespace), 1),
	}

	podList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	l := make([]v1.Pod, 0)
	for _, pi := range podList.Items {
		if pi.Name == podName {
			l = append(l, pi)
		}
	}

	events := filterEventsByPodsUID(eventList.Items, l)
	return FillEventsType(events), nil
}

// GetNodeEvents gets events associated to node with given name.
func GetNodeEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, nodeName string) (*common.EventList, error) {
	eventList := common.EventList{
		Events: make([]common.Event, 0),
	}

	scheme := runtime.NewScheme()
	groupVersion := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.AddKnownTypes(groupVersion, &v1.Node{})

	mc := client.CoreV1().Nodes()
	node, err := mc.Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err != nil {
		return &eventList, err
	}

	events, err := client.CoreV1().Events(v1.NamespaceAll).Search(scheme, node)
	_, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return &eventList, criticalError
	}

	eventList = CreateEventList(FillEventsType(events.Items), dsQuery)
	return &eventList, nil
}

// GetNamespaceEvents gets events associated to a namespace with given name.
func GetNamespaceEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace string) (common.EventList, error) {
	events, _ := client.CoreV1().Events(namespace).List(context.TODO(), api.ListEverything)
	return CreateEventList(FillEventsType(events.Items), dsQuery), nil
}

// FillEventsType is based on event Reason fills event Type in order to allow correct filtering by Type.
func FillEventsType(events []v1.Event) []v1.Event {
	for i := range events {
		// Fill in only events with empty type.
		if len(events[i].Type) == 0 {
			if isFailedReason(events[i].Reason, FailedReasonPartials...) {
				events[i].Type = v1.EventTypeWarning
			} else {
				events[i].Type = v1.EventTypeNormal
			}
		}
	}

	return events
}

// ToEvent converts event api Event to Event model object.
func ToEvent(event v1.Event) common.Event {
	result := common.Event{
		ObjectMeta:      api.NewObjectMeta(event.ObjectMeta),
		TypeMeta:        api.NewTypeMeta(api.ResourceKindEvent),
		Message:         event.Message,
		SourceComponent: event.Source.Component,
		SourceHost:      event.Source.Host,
		SubObject:       event.InvolvedObject.FieldPath,
		Count:           event.Count,
		FirstSeen:       event.FirstTimestamp,
		LastSeen:        event.LastTimestamp,
		Reason:          event.Reason,
		Type:            event.Type,
	}

	return result
}

// GetResourceEvents gets events associated to specified resource.
func GetResourceEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace, name string) (*common.EventList, error) {
	resourceEvents, err := GetEvents(client, namespace, name)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return EmptyEventList, err
	}

	events := CreateEventList(resourceEvents, dsQuery)
	events.Errors = nonCriticalErrors
	return &events, nil
}

// CreateEventList converts array of api events to common EventList structure
func CreateEventList(events []v1.Event, dsQuery *dataselect.DataSelectQuery) common.EventList {
	eventList := common.EventList{
		Events:   make([]common.Event, 0),
		ListMeta: api.ListMeta{TotalItems: len(events)},
	}

	events = fromCells(dataselect.GenericDataSelect(toCells(events), dsQuery))
	for _, event := range events {
		eventDetail := ToEvent(event)
		eventList.Events = append(eventList.Events, eventDetail)
	}

	return eventList
}

// The code blow allows to perform complex data section on []api.Event
type EventCell v1.Event

func (self EventCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		// if name is not supported then just return a constant dummy value, sort will have no effect.
		return nil
	}
}

func toCells(std []v1.Event) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = EventCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.Event {
	std := make([]v1.Event, len(cells))
	for i := range std {
		std[i] = v1.Event(cells[i].(EventCell))
	}
	return std
}
