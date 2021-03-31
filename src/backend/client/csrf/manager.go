package csrf

import (
	"context"
	"github.com/cuijxin/k8s-dashboard/src/backend/args"
	"github.com/cuijxin/k8s-dashboard/src/backend/client/api"

	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Implements CsrfTokenManager interface.
type csrfTokenManager struct {
	token  string
	client kubernetes.Interface
}

func (c *csrfTokenManager) init() {
	log.Printf("Initializing csrf token from %s secret", api.CsrfTokenSecretName)
	tokenSecret, err := c.client.CoreV1().Secrets(args.Holder.GetNamespace()).
		Get(context.TODO(), api.CsrfTokenSecretName, v1.GetOptions{})

	if err != nil {
		panic(err)
	}

	token := string(tokenSecret.Data[api.CsrfTokenSecretData])
	if len(token) == 0 {
		log.Printf("Empty token. Generating and storing in a secret %s", api.CsrfTokenSecretName)
		token = api.GenerateCSRFKey()
		tokenSecret.StringData = map[string]string{api.CsrfTokenSecretData: token}
		_, err := c.client.CoreV1().Secrets(args.Holder.GetNamespace()).Update(context.TODO(), tokenSecret, v1.UpdateOptions{})
		if err != nil {
			panic(err)
		}
	}

	c.token = token
}

// Token implements CsrfTokenManager interface.
func (c *csrfTokenManager) Token() string {
	return c.token
}

// NewCsrfTokenManager creates and initializes new instance of csrf token manager.
func NewCsrfTokenManager(client kubernetes.Interface) api.CsrfTokenManager {
	manager := &csrfTokenManager{client: client}
	manager.init()

	return manager
}
