package feature

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Feature = Feature{Name: "mark", Title: "title", Description: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name    string
		feature Feature
		err     string
		code    int
	}{
		{name: "ok", feature: u, err: ""},
		{name: "Empty Name", feature: Feature{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty title", feature: Feature{Name: "feature name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Description", feature: Feature{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.feature.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
