package location

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatelocation(t *testing.T) {
	location := &Location{}
	if err := json.Unmarshal([]byte(jsondata), &location); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewlocationService(NewlocationRepo())
	u, err := service.Create(location)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAlllocation(t *testing.T) {

	location := &Location{}
	if err := json.Unmarshal([]byte(jsondata), &location); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewlocationService(NewlocationRepo())
	location1, _ := service.Create(location)
	_, err := service.GetAll()
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(location1.Code)
}
func TestServiceGetOnelocation(t *testing.T) {

	location := &Location{}
	if err := json.Unmarshal([]byte(jsondata), &location); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewlocationService(NewlocationRepo())
	location1, err := service.Create(location)
	assert.EqualValues(t, "name", location.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(location1.Code)
}

func TestServiceUpdatelocation(t *testing.T) {

	location := &Location{}
	if err := json.Unmarshal([]byte(jsondata), &location); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewlocationService(NewlocationRepo())
	location1, err := service.Create(location)
	assert.Nil(t, err)
	location1.Name = "updatedname"
	res, err := service.Update(location1.Code, location1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(location1.Code)
}
func TestServiceDeletelocation(t *testing.T) {

	location := &Location{}
	if err := json.Unmarshal([]byte(jsondata), &location); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewlocationService(NewlocationRepo())
	location1, err := service.Create(location)
	res, err := service.Delete(location1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
