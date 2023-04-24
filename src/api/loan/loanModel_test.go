package loan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Loan = Loan{Name: "mark", Business: "http://example.com", Amount: 30000}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		loan Loan
		err  string
		code int
	}{
		{name: "ok", loan: u, err: ""},
		{name: "Empty Name", loan: Loan{Name: "", Business: "asd455", Amount: 30000}, err: "Name should not be empty", code: 400},
		{name: "Empty loan pin", loan: Loan{Name: "loan name", Business: "", Amount: 30000}, err: "Business Name should not be empty", code: 400},
		{name: "Empty Address", loan: Loan{Name: "name", Business: "white", Amount: 0}, err: "Amount should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.loan.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
