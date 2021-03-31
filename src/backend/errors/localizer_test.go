package errors

import "testing"

func TestLocalizeError(t *testing.T) {
	cases := []struct {
		err      error
		expected error
	}{
		{
			nil,
			nil,
		},
		{
			NewInternal("some unknown error"),
			NewInternal("some unknown error"),
		},
		{
			NewInvalid("does not match the namespace"),
			NewInvalid("MSG_DEPLOY_NAMESPACE_MISMATCH_ERROR"),
		},
		{
			NewInvalid("empty namespace may not be set"),
			NewInvalid("MSG_DEPLOY_EMPTY_NAMESPACE_ERROR"),
		},
	}
	for _, c := range cases {
		actual := LocalizeError(c.err)
		if !areErrorsEqual(actual, c.expected) {
			t.Errorf("LocalizeError(%+v) == %+v, expected %+v", c.err, actual, c.expected)
		}
	}
}

func areErrorsEqual(err1, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() == err2.Error()) ||
		(err1 == nil && err2 == nil)
}
