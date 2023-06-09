package news

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var u News = News{Name: "mark", Title: "http://example.com", Content: "addr456"}

func TestValidateUserInputRequiredFields(t *testing.T) {
	testcases := []struct {
		name string
		news News
		err  string
		code int
	}{
		{name: "ok", news: u, err: ""},
		{name: "Empty Name", news: News{Name: "", Content: "asd455", Title: "http://example.com"}, err: "Name should not be empty", code: 400},
		{name: "Empty news pin", news: News{Name: "news name", Content: "", Title: "http://example.com"}, err: "Business Pin should not be empty", code: 400},
		{name: "Empty Address", news: News{Name: "name", Content: "white", Title: ""}, err: "Url endpoint should not be empty", code: 400},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := test.news.Validate()
			if err != nil {
				require.EqualValues(t, test.err, err.Message())
				require.EqualValues(t, test.code, err.Code())
			}
		})
	}

}
