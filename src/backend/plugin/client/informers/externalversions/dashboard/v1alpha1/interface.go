// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/cuijxin/k8s-dashboard/src/backend/plugin/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Plugins returns a PluginInformer.
	Plugins() PluginInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Plugins returns a PluginInformer.
func (v *version) Plugins() PluginInformer {
	return &pluginInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
