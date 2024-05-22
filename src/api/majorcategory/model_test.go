package majorcategory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMajorcategoryInputRequiredFields(t *testing.T) {
	jsondata := `{"name":"name","title":"title","description":"desscription","shopalias": "alias"}`
	majorcategory := &Majorcategory{}
	if err := json.Unmarshal([]byte(jsondata), &majorcategory); err != nil {
		t.Errorf("failed to unmarshal tag data %v", err.Error())
	}
	// fmt.Println("------------------", tag)
	expected := ""
	if err := majorcategory.Validate(); err != nil {
		expected = "Invalid Name"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Erro valivating Name")
		}
		expected = "Invalid title"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating title")
		}
		expected = "Invalid Description"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Description")
		}
		expected = "Invalid Shopalias"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Shopalias")
		}
	}
}
