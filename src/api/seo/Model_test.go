package seo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Seo = Seo{Name: "mark", Title: "title", Description: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		seo  Seo
		err  string
		code int
	}{
		{name: "ok", seo: u, err: ""},
		{name: "Empty Name", seo: Seo{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty title", seo: Seo{Name: "seo name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Description", seo: Seo{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.seo.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
