package rating

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"productname":"name","productcode":"code123","description":"description", "rate":4}`

func TestServiceCreateRating(t *testing.T) {
	// fmt.Println(">>>>>>>>>", Bizname)
	// Bizname = "test"
	rating := &Rating{}
	if err := json.Unmarshal([]byte(jsondata), &rating); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewratingService(NewratingRepo())
	Bizname = "alias4567"
	// fmt.Println(">>>>>>>>>rating", *rating)
	u, err := service.Create(rating)
	// fmt.Println(">>>>>>>>>", u)
	assert.EqualValues(t, "name", u.Productname, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllrating(t *testing.T) {

	rating := &Rating{}
	if err := json.Unmarshal([]byte(jsondata), &rating); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewratingService(NewratingRepo())
	Bizname = "alias4567"
	_, _ = service.Create(rating)
	_, err := service.GetAll()
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(rating.Code)
}
func TestServiceGetOnerating(t *testing.T) {

	rating := &Rating{}
	if err := json.Unmarshal([]byte(jsondata), &rating); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewratingService(NewratingRepo())
	Bizname = "alias4567"
	rating, err := service.Create(rating)
	assert.EqualValues(t, "name", rating.Productname, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(rating.Code)
}

func TestServiceUpdaterating(t *testing.T) {
	rating := &Rating{}
	if err := json.Unmarshal([]byte(jsondata), &rating); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewratingService(NewratingRepo())
	Bizname = "alias4567"
	rating1, _ := service.Create(rating)
	rating1.Productname = "update"
	u, err := service.Update(rating.Code, rating1)
	assert.EqualValues(t, rating1.Productname, u.Productname, "Something went wrong testing with the Getting one method")
	assert.Nil(t, err)
	service.Delete(rating1.Code)
}
func TestServiceDeleterating(t *testing.T) {

	rating1 := &Rating{}
	if err := json.Unmarshal([]byte(jsondata), &rating1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewratingService(NewratingRepo())
	Bizname = "alias4567"
	rating, _ := service.Create(rating1)
	res, err := service.Delete(rating.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
