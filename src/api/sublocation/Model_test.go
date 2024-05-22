package subLocation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u SubLocation = SubLocation{Name: "mark", Title: "title", Description: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name        string
		subLocation SubLocation
		err         string
		code        int
	}{
		{name: "ok", subLocation: u, err: ""},
		{name: "Empty Name", subLocation: SubLocation{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty title", subLocation: SubLocation{Name: "subLocation name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Description", subLocation: SubLocation{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.subLocation.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
