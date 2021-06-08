package customresourcedefinition

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/common"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/resource/event"
	"github.com/yanniszark/go-nodetool/client"
)

// GetEventsForCustomResourceObject gets events that are associated with this CR object.
func GetEventsForCustomResourceObject(client client.Interface,
	dsQuery *dataselect.DataSelectQuery,
	namespace, name string) (*common.EventList, error) {
	return event.GetResourceEvents(client, dsQuery, namespace, name)
}
