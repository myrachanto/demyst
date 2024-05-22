package blog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u Blog = Blog{Name: "mark", Title: "http://example.com", Content: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		Blog Blog
		err  string
		code int
	}{
		{name: "ok", Blog: u, err: ""},
		{name: "Empty Name", Blog: Blog{Name: "", Content: "asd455", Title: "http://example.com"}, err: "Name should not be empty", code: 400},
		{name: "Empty Blog pin", Blog: Blog{Name: "Blog name", Content: "", Title: "http://example.com"}, err: "Business Pin should not be empty", code: 400},
		{name: "Empty Address", Blog: Blog{Name: "name", Content: "white", Title: ""}, err: "Url endpoint should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.Blog.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
