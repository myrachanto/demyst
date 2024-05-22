package gift

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/estate/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreateGift(t *testing.T) {
	Gift := &Gift{}
	if err := json.Unmarshal([]byte(jsondata), &Gift); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewGiftService(NewGiftRepo())
	u, err := service.Create(Gift)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllGift(t *testing.T) {

	Gift := &Gift{}
	if err := json.Unmarshal([]byte(jsondata), &Gift); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewGiftService(NewGiftRepo())
	Gift1, _ := service.Create(Gift)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(Gift1.Code)
}
func TestServiceGetOneGift(t *testing.T) {

	Gift := &Gift{}
	if err := json.Unmarshal([]byte(jsondata), &Gift); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewGiftService(NewGiftRepo())
	Gift1, err := service.Create(Gift)
	assert.EqualValues(t, "name", Gift.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(Gift1.Code)
}

func TestServiceUpdateGift(t *testing.T) {

	Gift := &Gift{}
	if err := json.Unmarshal([]byte(jsondata), &Gift); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewGiftService(NewGiftRepo())
	Gift1, err := service.Create(Gift)
	assert.Nil(t, err)
	Gift1.Name = "updatedname"
	res, err := service.Update(Gift1.Code, Gift1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(Gift1.Code)
}
func TestServiceDeleteGift(t *testing.T) {

	Gift := &Gift{}
	if err := json.Unmarshal([]byte(jsondata), &Gift); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewGiftService(NewGiftRepo())
	Gift1, err := service.Create(Gift)
	res, err := service.Delete(Gift1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
