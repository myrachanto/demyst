package gift

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Gift = Gift{Name: "mark", Title: "title", Description: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name     string
		Gift Gift
		err      string
		code     int
	}{
		{name: "ok", Gift: u, err: ""},
		{name: "Empty Name", Gift: Gift{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty title", Gift: Gift{Name: "Gift name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Description", Gift: Gift{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.Gift.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
