package settings

import (
	"reflect"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/settings/api"

	"k8s.io/client-go/kubernetes/fake"
)

func TestNewSettingsManager(t *testing.T) {
	sm := NewSettingsManager().(*SettingsManager)

	if len(sm.settings) > 0 {
		t.Error("new settigns manager should have no settings set")
	}
}

func TestSettingsManager_SaveGlobalSettings(t *testing.T) {
	sm := NewSettingsManager()
	client := fake.NewSimpleClientset(api.GetDefaultSettingsConfigMap(""))
	gs := sm.GetGlobalSettings(client)

	if !reflect.DeepEqual(api.GetDefaultSettings(), gs) {
		t.Errorf("it should return default settings \"%v\" instead of \"%v\"", api.GetDefaultSettings(), gs)
	}
}

func TestSettignsManager_SaveGlobalSettings(t *testing.T) {
	sm := NewSettingsManager()
	client := fake.NewSimpleClientset(api.GetDefaultSettingsConfigMap(""))
	defaults := api.GetDefaultSettings()
	err := sm.SaveGlobalSettings(client, &defaults)

	if err == nil {
		t.Errorf("it should fail with \"%s\" error if trying to save but manager has deprecated data",
			api.ConcurrentSettingsChangeError)
	}

	if !reflect.DeepEqual(err.Error(), api.ConcurrentSettingsChangeError) {
		t.Errorf("it should fail with \"%s\" error instead of \"%s\" if trying to save but manager has deprecated data",
			api.ConcurrentSettingsChangeError, err.Error())
	}

	err = sm.SaveGlobalSettings(client, &defaults)

	if err != nil {
		t.Errorf("it should save settings if manager has no deprecated data instead of failing with \"%s\", error",
			err.Error())
	}
}
