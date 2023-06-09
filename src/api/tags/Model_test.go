package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Tag = Tag{Name: "name", Title: "title", Description: "description"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		tag  Tag
		err  string
		code int
	}{
		{name: "ok", tag: u, err: ""},
		{name: "Empty Name", tag: Tag{Name: "", Description: "asd455", Title: "title"}, err: "Name should not be empty", code: 400},
		{name: "Empty description", tag: Tag{Name: "tag name", Description: "", Title: "title"}, err: "Description should not be empty", code: 400},
		{name: "Empty Title", tag: Tag{Name: "name", Description: "white", Title: ""}, err: "Title should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.tag.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
