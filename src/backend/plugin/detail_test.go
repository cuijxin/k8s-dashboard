package plugin

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/plugin/apis/dashboard/v1alpha1"
	fakePluginClientset "github.com/cuijxin/k8s-dashboard/src/backend/plugin/client/clientset/versioned/fake"
	"github.com/emicklei/go-restful/v3"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeK8sClient "k8s.io/client-go/kubernetes/fake"
)

var srcData = "randomPluginSourceCode"

func TestGetPluginSource(t *testing.T) {
	ns := "default"
	pluginName := "test-plugin"
	filename := "plugin-test.js"
	cfgMapName := "plugin-test-cfgMap"

	pcs := fakePluginClientset.NewSimpleClientset()
	cs := fakeK8sClient.NewSimpleClientset()

	_, err := GetPluginSource(pcs, cs, ns, pluginName)
	if err == nil {
		t.Errorf("error 'plugins.dashboard.k8s.io \"%s\" not found' did not occur", pluginName)
	}

	_, _ = pcs.DashboardV1alpha1().Plugins(ns).Create(context.TODO(), &v1alpha1.Plugin{
		ObjectMeta: v1.ObjectMeta{Name: pluginName, Namespace: ns},
		Spec: v1alpha1.PluginSpec{
			Source: v1alpha1.Source{
				ConfigMapRef: &coreV1.ConfigMapEnvSource{
					LocalObjectReference: coreV1.LocalObjectReference{Name: cfgMapName},
				},
				Filename: filename}},
	}, v1.CreateOptions{})

	_, err = GetPluginSource(pcs, cs, ns, pluginName)
	if err == nil {
		t.Errorf("error 'configmaps \"%s\" not found' did not occur", cfgMapName)
	}

	_, _ = cs.CoreV1().ConfigMaps(ns).Create(context.TODO(), &coreV1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: cfgMapName, Namespace: ns},
		Data: map[string]string{filename: srcData},
	}, v1.CreateOptions{})

	data, err := GetPluginSource(pcs, cs, ns, pluginName)
	if err != nil {
		t.Errorf("error while fetching plugin source: %s", err)
	}

	if !bytes.Equal(data, []byte(srcData)) {
		t.Error("bytes in configMap and bytes from GetPluginSource are different")
	}
}

func Test_servePluginSource(t *testing.T) {
	ns := "default"
	pluginName := "test-plugin"
	filename := "plugin-test.js"
	cfgMapName := "plugin-test-cfgMap"
	h := Handler{&fakeClientManager{}}

	pcs, _ := h.cManager.PluginClient(nil)
	_, _ = pcs.DashboardV1alpha1().Plugins(ns).Create(context.TODO(), &v1alpha1.Plugin{
		ObjectMeta: v1.ObjectMeta{Name: pluginName, Namespace: ns},
		Spec: v1alpha1.PluginSpec{
			Source: v1alpha1.Source{
				ConfigMapRef: &coreV1.ConfigMapEnvSource{
					LocalObjectReference: coreV1.LocalObjectReference{Name: cfgMapName},
				},
				Filename: filename}},
	}, v1.CreateOptions{})

	httpReq, _ := http.NewRequest(http.MethodGet, "/api/v1/plugin/default/test-plugin", nil)
	req := restful.NewRequest(httpReq)

	httpWriter := httptest.NewRecorder()
	resp := restful.NewResponse(httpWriter)

	h.servePluginSource(req, resp)
}
