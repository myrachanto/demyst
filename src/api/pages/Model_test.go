package pages

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Page = Page{Name: "name", Title: "title", Content: "content"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		page Page
		err  string
		code int
	}{
		{name: "ok", page: u, err: ""},
		{name: "Empty Name", page: Page{Name: "", Content: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty page pin", page: Page{Name: "page name", Content: "", Title: "title"}, err: "Title should not be empty", code: 400},
		{name: "Empty Address", page: Page{Name: "name", Content: "white", Title: ""}, err: "Content should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.page.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
