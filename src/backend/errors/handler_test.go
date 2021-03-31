package errors

import (
	"reflect"
	"testing"
)

func TestHandleHTTPError(t *testing.T) {
	cases := []struct {
		err      error
		expected int
	}{
		{
			nil,
			500,
		},
		{
			NewInvalid("some unknown error"),
			500,
		},
		{
			NewInvalid(MsgLoginUnauthorizedError),
			500,
		},
		{
			NewInvalid(MsgTokenExpiredError),
			401,
		},
		{
			NewInvalid(MsgEncryptionKeyChanged),
			401,
		},
	}

	for _, c := range cases {
		actual := HandleHTTPError(c.err)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("HandleHTTPError(%+v) == %+v, expected %+v", c.err, actual, c.expected)
		}
	}
}
