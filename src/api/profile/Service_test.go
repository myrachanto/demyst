package profile

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/estate/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreateProfile(t *testing.T) {
	Profile := &Profile{}
	if err := json.Unmarshal([]byte(jsondata), &Profile); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewprofileService(NewProfileRepo())
	u, err := service.Create(Profile)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllProfile(t *testing.T) {

	Profile := &Profile{}
	if err := json.Unmarshal([]byte(jsondata), &Profile); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewprofileService(NewProfileRepo())
	Profile1, _ := service.Create(Profile)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(Profile1.Code)
}
func TestServiceGetOneProfile(t *testing.T) {

	Profile := &Profile{}
	if err := json.Unmarshal([]byte(jsondata), &Profile); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewprofileService(NewProfileRepo())
	Profile1, err := service.Create(Profile)
	assert.EqualValues(t, "name", Profile.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(Profile1.Code)
}

func TestServiceUpdateProfile(t *testing.T) {

	Profile := &Profile{}
	if err := json.Unmarshal([]byte(jsondata), &Profile); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewprofileService(NewProfileRepo())
	Profile1, err := service.Create(Profile)
	assert.Nil(t, err)
	Profile1.Name = "updatedname"
	res, err := service.Update(Profile1.Code, Profile1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(Profile1.Code)
}
func TestServiceDeleteProfile(t *testing.T) {

	Profile := &Profile{}
	if err := json.Unmarshal([]byte(jsondata), &Profile); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewprofileService(NewProfileRepo())
	Profile1, err := service.Create(Profile)
	res, err := service.Delete(Profile1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
