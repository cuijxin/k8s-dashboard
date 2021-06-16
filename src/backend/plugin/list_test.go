package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/plugin/apis/dashboard/v1alpha1"
	fakePluginClientset "github.com/cuijxin/k8s-dashboard/src/backend/plugin/client/clientset/versioned/fake"
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
	"github.com/emicklei/go-restful/v3"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPluginList(t *testing.T) {
	ns := "default"
	pluginName := "test-plugin"
	filename := "plugin-test.js"
	cfgMapName := "plugin-test-cfgMap"

	pcs := fakePluginClientset.NewSimpleClientset()

	_, _ = pcs.DashboardV1alpha1().Plugins(ns).Create(context.TODO(), &v1alpha1.Plugin{
		ObjectMeta: v1.ObjectMeta{Name: pluginName, Namespace: ns},
		Spec: v1alpha1.PluginSpec{
			Source: v1alpha1.Source{
				ConfigMapRef: &coreV1.ConfigMapEnvSource{
					LocalObjectReference: coreV1.LocalObjectReference{Name: cfgMapName},
				},
				Filename: filename}},
	}, metaV1.CreateOptions{})

	dsQuery := dataselect.DataSelectQuery{
		PaginationQuery: &dataselect.PaginationQuery{
			ItemsPerPage: 10,
			Page:         1,
		},
		SortQuery:   dataselect.NoSort,
		FilterQuery: dataselect.NoFilter,
		MetricQuery: dataselect.NoMetrics,
	}
	data, err := GetPluginList(pcs, ns, &dsQuery)
	if err != nil {
		t.Errorf("error while fetching plugins: %s", err)
	}

	if data.ListMeta.TotalItems != 1 {
		t.Errorf("there should be one plugin registered, got %d", data.ListMeta.TotalItems)
	}
}

func Test_handlePluginList(t *testing.T) {
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
	}, metaV1.CreateOptions{})

	httpReq, _ := http.NewRequest(http.MethodGet, "/api/v1/plugin/default?itemsPerPage=10&page=1&sortBy=d,creationTimestamp", nil)
	req := restful.NewRequest(httpReq)

	httpWriter := httptest.NewRecorder()
	resp := restful.NewResponse(httpWriter)

	h.handlePluginList(req, resp)
}
