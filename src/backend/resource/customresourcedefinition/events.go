package customresourcedefinition

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/event"
	client "k8s.io/client-go/kubernetes"
)

// GetEventsForCustomResourceObject gets events that are associated with this CR object.
func GetEventsForCustomResourceObject(client client.Interface,
	dsQuery *dataselect.DataSelectQuery,
	namespace, name string) (*common.EventList, error) {
	return event.GetResourceEvents(client, dsQuery, namespace, name)
}
