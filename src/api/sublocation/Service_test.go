package subLocation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatesubLocation(t *testing.T) {
	subLocation := &SubLocation{}
	if err := json.Unmarshal([]byte(jsondata), &subLocation); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewsubLocationService(NewsubLocationRepo())
	u, err := service.Create(subLocation)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllsubLocation(t *testing.T) {

	subLocation := &SubLocation{}
	if err := json.Unmarshal([]byte(jsondata), &subLocation); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewsubLocationService(NewsubLocationRepo())
	subLocation1, _ := service.Create(subLocation)
	_, err := service.GetAll()
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(subLocation1.Code)
}
func TestServiceGetOnesubLocation(t *testing.T) {

	subLocation := &SubLocation{}
	if err := json.Unmarshal([]byte(jsondata), &subLocation); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewsubLocationService(NewsubLocationRepo())
	subLocation1, err := service.Create(subLocation)
	assert.EqualValues(t, "name", subLocation.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(subLocation1.Code)
}

func TestServiceUpdatesubLocation(t *testing.T) {

	subLocation := &SubLocation{}
	if err := json.Unmarshal([]byte(jsondata), &subLocation); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewsubLocationService(NewsubLocationRepo())
	subLocation1, err := service.Create(subLocation)
	assert.Nil(t, err)
	subLocation1.Name = "updatedname"
	res, err := service.Update(subLocation1.Code, subLocation1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(subLocation1.Code)
}
func TestServiceDeletesubLocation(t *testing.T) {

	subLocation := &SubLocation{}
	if err := json.Unmarshal([]byte(jsondata), &subLocation); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewsubLocationService(NewsubLocationRepo())
	subLocation1, err := service.Create(subLocation)
	res, err := service.Delete(subLocation1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
