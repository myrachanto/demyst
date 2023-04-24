package accounting

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Accounting = Accounting{Name: "mark", UrlEndpoint: "http://example.com", BusinessPin: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name       string
		accounting Accounting
		err        string
		code       int
	}{
		{name: "ok", accounting: u, err: ""},
		{name: "Empty Name", accounting: Accounting{Name: "", BusinessPin: "asd455", UrlEndpoint: "http://example.com"}, err: "Name should not be empty", code: 400},
		{name: "Empty accounting pin", accounting: Accounting{Name: "accounting name", BusinessPin: "", UrlEndpoint: "http://example.com"}, err: "Business Pin should not be empty", code: 400},
		{name: "Empty Address", accounting: Accounting{Name: "name", BusinessPin: "white", UrlEndpoint: ""}, err: "Url endpoint should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.accounting.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
