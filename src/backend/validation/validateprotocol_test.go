package validation

import "testing"

func TestValidateProtocol(t *testing.T) {
	cases := []struct {
		spec     *ProtocolValiditySpec
		expected bool
	}{
		{
			&ProtocolValiditySpec{
				Protocol:   "TCP",
				IsExternal: false,
			},
			true,
		},
		{
			&ProtocolValiditySpec{
				Protocol:   "UDP",
				IsExternal: true,
			},
			false,
		},
	}

	for _, c := range cases {
		validity := ValidateProtocol(c.spec)
		if validity.Valid != c.expected {
			t.Errorf("Expected %#v validity to be %#v, but was %#v\n", c.spec, c.expected, validity)
		}
	}
}
