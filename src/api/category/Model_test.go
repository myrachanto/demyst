package category

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Category = Category{Name: "mark", Title: "title", Description: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name     string
		category Category
		err      string
		code     int
	}{
		{name: "ok", category: u, err: ""},
		{name: "Empty Name", category: Category{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty title", category: Category{Name: "category name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Description", category: Category{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.category.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
