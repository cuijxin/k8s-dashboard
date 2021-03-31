package csrf

import (
	"testing"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/fake"

	"github.com/cuijxin/k8s-dashboard/src/backend/client/api"
)

func TestCsrfTokenManager_Token(t *testing.T) {
	cases := []struct {
		info       string
		csrfSecret *v1.Secret
		wantPanic  bool
		wantToken  bool
	}{
		{"should panic when secret does not exist", nil, true, false},
		{"should generate token when secret exists",
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: api.CsrfTokenSecretName,
				},
			}, false, true},
	}

	for _, c := range cases {
		t.Run(c.info, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != c.wantPanic {
					t.Errorf("Recover = %v, wantPanic = %v", r, c.wantPanic)
				}
			}()

			client := fake.NewSimpleClientset(c.csrfSecret)
			manager := NewCsrfTokenManager(client)

			if (len(manager.Token()) == 0) == c.wantPanic {
				t.Errorf("Expected token to exists: %v", c.wantToken)
			}
		})
	}
}
