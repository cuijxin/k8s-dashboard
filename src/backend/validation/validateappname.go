package validation

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"

	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
)

// AppNameValiditySpec is a specification for application name validation request.
type AppNameValiditySpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// AppNameValidity describes validity of the applicaiton name.
type AppNameValidity struct {
	// True when the applicaiton name is valid.
	Valid bool `json:"valid"`
}

// ValidateAppName validates application name. When error is returned, name
// validity could not be determined.
func ValidateAppName(spec *AppNameValiditySpec, client client.Interface) (*AppNameValidity, error) {
	log.Printf("Validating %s application name in %s namespace", spec.Name, spec.Namespace)

	isValidDeployment := false
	isValidService := false

	_, err := client.AppsV1().Deployments(spec.Namespace).Get(context.TODO(), spec.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFoundError(err) || errors.IsForbiddenError(err) {
			isValidDeployment = true
		} else {
			return nil, err
		}
	}

	_, err = client.CoreV1().Services(spec.Namespace).Get(context.TODO(), spec.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFoundError(err) || errors.IsForbiddenError(err) {
			isValidService = true
		} else {
			return nil, err
		}
	}

	isValid := isValidDeployment && isValidService

	log.Printf("Validation result for %s applicaiton name in %s namespace is %t", spec.Name, spec.Namespace, isValid)

	return &AppNameValidity{Valid: isValid}, nil
}
