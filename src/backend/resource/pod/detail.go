package pod

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/api"
	"github.com/kubernetes/dashboard/src/app/backend/resource/controller"
)

// PodDetail is a presentation layer view of Kubernetes Pod resource.
type PodDetail struct {
	ObjectMeta         api.ObjectMeta            `json:"objectMeta"`
	TypeMeta           api.TypeMeta              `json:"typeMeta"`
	PodPhase           string                    `json:"podPhase"`
	PodIP              string                    `json:"podIP"`
	NodeName           string                    `json:"nodeName"`
	ServiceAccountName string                    `json:"serviceAccountName"`
	RestartCount       int32                     `json:"restartCount"`
	QOSClass           string                    `json:"qosClass"`
	Controller         *controller.ResourceOwner `json:"controller,omitempty"`
}
