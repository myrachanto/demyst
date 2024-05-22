package feature

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatefeature(t *testing.T) {
	feature := &Feature{}
	if err := json.Unmarshal([]byte(jsondata), &feature); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewfeatureService(NewfeatureRepo())
	u, err := service.Create(feature)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllfeature(t *testing.T) {

	feature := &Feature{}
	if err := json.Unmarshal([]byte(jsondata), &feature); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewfeatureService(NewfeatureRepo())
	feature1, _ := service.Create(feature)
	_, err := service.GetAll("")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(feature1.Code)
}
func TestServiceGetOnefeature(t *testing.T) {

	feature := &Feature{}
	if err := json.Unmarshal([]byte(jsondata), &feature); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewfeatureService(NewfeatureRepo())
	feature1, err := service.Create(feature)
	assert.EqualValues(t, "name", feature.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(feature1.Code)
}

func TestServiceUpdatefeature(t *testing.T) {

	feature := &Feature{}
	if err := json.Unmarshal([]byte(jsondata), &feature); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewfeatureService(NewfeatureRepo())
	feature1, err := service.Create(feature)
	assert.Nil(t, err)
	feature1.Name = "updatedname"
	res, err := service.Update(feature1.Code, feature1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(feature1.Code)
}
func TestServiceDeletefeature(t *testing.T) {

	feature := &Feature{}
	if err := json.Unmarshal([]byte(jsondata), &feature); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewfeatureService(NewfeatureRepo())
	feature1, err := service.Create(feature)
	res, err := service.Delete(feature1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
