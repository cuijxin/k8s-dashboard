package validation

import (
	"net/http"
	"net/textproto"
	"reflect"
	"testing"

	"github.com/cuijxin/k8s-dashboard/src/backend/client"
	restful "github.com/emicklei/go-restful/v3"
)

func TestValidateLoginStatus(t *testing.T) {
	cases := []struct {
		info     string
		request  *restful.Request
		expected *LoginStatus
	}{
		{
			"Should indicate that user is logged in with token",
			&restful.Request{Request: &http.Request{Header: http.Header(map[string][]string{
				textproto.CanonicalMIMEHeaderKey(client.JWETokenHeader): {"test-token"},
			})}},
			&LoginStatus{TokenPresent: true},
		},
	}

	for _, c := range cases {
		status := ValidateLoginStatus(c.request)

		if !reflect.DeepEqual(status, c.expected) {
			t.Errorf("Test Case: %s. Expected status to be: %v, but got %v.", c.info, c.expected, status)
		}
	}
}
