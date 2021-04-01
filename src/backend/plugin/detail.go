package plugin

import (
	"context"

	pluginclientset "github.com/cuijxin/k8s-dashboard/src/backend/plugin/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPluginSource has the logic to get the actual plugin source code from information in Plugin.Spec
func GetPluginSource(client pluginclientset.Interface, k8sClient kubernetes.Interface, ns string, name string) ([]byte, error) {
	plugin, err := client.DashboardV1alpha1().Plugins(ns).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	cfgMap, err := k8sClient.CoreV1().ConfigMaps(ns).Get(context.TODO(), plugin.Spec.Source.ConfigMapRef.Name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return []byte(cfgMap.Data[plugin.Spec.Source.Filename]), nil
}
