package rating

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRatingInputRequiredFields(t *testing.T) {
	jsondata := `{"productcode":"code1234","rate":3,"description":"desscription","shopalias": "alias"
	}`
	rating := &Rating{}
	if err := json.Unmarshal([]byte(jsondata), &rating); err != nil {
		t.Errorf("failed to unmarshal tag data %v", err.Error())
	}
	// fmt.Println("------------------", tag)
	expected := ""
	if err := rating.Validate(); err != nil {
		expected = "Invalid Code"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Erro valivating Code")
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
