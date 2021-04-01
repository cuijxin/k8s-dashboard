package validation

import (
	"github.com/docker/distribution/reference"
)

// ImageReferenceValiditySpec is a specification of an image reference validation request.
type ImageReferenceValiditySpec struct {
	// Reference of the image
	Reference string `json:"reference"`
}

// ImageReferenceValidity describes validity of the image reference.
type ImageReferenceValidity struct {
	// True when the image reference is valid.
	Valid bool `json:"valid"`
	// Error reason when image reference is valid
	Reason string `json:"reason"`
}

// ValidateImageReference validates image reference.
func ValidateImageReference(spec *ImageReferenceValiditySpec) (*ImageReferenceValidity, error) {
	s := spec.Reference
	_, err := reference.Parse(s)
	if err != nil {
		return &ImageReferenceValidity{Valid: false, Reason: err.Error()}, nil
	}
	return &ImageReferenceValidity{Valid: true}, nil
}
