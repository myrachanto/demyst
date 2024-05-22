package profile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Profile = Profile{Name: "mark", Title: "http://example.com", Content: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		Profile Profile
		err  string
		code int
	}{
		{name: "ok", Profile: u, err: ""},
		{name: "Empty Name", Profile: Profile{Name: "", Content: "asd455", Title: "http://example.com"}, err: "Name should not be empty", code: 400},
		{name: "Empty Profile pin", Profile: Profile{Name: "Profile name", Content: "", Title: "http://example.com"}, err: "Business Pin should not be empty", code: 400},
		{name: "Empty Address", Profile: Profile{Name: "name", Content: "white", Title: ""}, err: "Url endpoint should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.Profile.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
