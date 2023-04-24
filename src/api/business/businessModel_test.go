package business

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Business = Business{Name: "mark", BusinessPin: "white", YearEstablished: 2020}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name     string
		business Business
		err      string
		code     int
	}{
		{name: "ok", business: u, err: ""},
		{name: "Empty Name", business: Business{Name: "", BusinessPin: "white", YearEstablished: 2020}, err: "Name should not be empty", code: 400},
		{name: "Empty Business pin", business: Business{Name: "business name", BusinessPin: "", YearEstablished: 2020}, err: "Business Pin should not be empty", code: 400},
		{name: "Empty Address", business: Business{Name: "name", BusinessPin: "white", YearEstablished: 0}, err: "Year of Establishment should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.business.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
