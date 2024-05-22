package location

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Location = Location{Name: "mark", Title: "title", Description: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name     string
		location Location
		err      string
		code     int
	}{
		{name: "ok", location: u, err: ""},
		{name: "Empty Name", location: Location{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty title", location: Location{Name: "location name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Description", location: Location{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.location.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
