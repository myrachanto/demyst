package seo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreateseo(t *testing.T) {
	seo := &Seo{}
	if err := json.Unmarshal([]byte(jsondata), &seo); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewseoService(NewseoRepo())
	u, err := service.Create(seo)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllseo(t *testing.T) {

	seo := &Seo{}
	if err := json.Unmarshal([]byte(jsondata), &seo); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewseoService(NewseoRepo())
	seo1, _ := service.Create(seo)
	_, err := service.GetAll("")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(seo1.Code)
}
func TestServiceGetOneseo(t *testing.T) {

	seo := &Seo{}
	if err := json.Unmarshal([]byte(jsondata), &seo); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewseoService(NewseoRepo())
	seo1, err := service.Create(seo)
	assert.EqualValues(t, "name", seo.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(seo1.Code)
}

func TestServiceUpdateseo(t *testing.T) {

	seo := &Seo{}
	if err := json.Unmarshal([]byte(jsondata), &seo); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewseoService(NewseoRepo())
	seo1, err := service.Create(seo)
	assert.Nil(t, err)
	seo1.Name = "updatedname"
	res, err := service.Update(seo1.Code, seo1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(seo1.Code)
}
func TestServiceDeleteseo(t *testing.T) {

	seo := &Seo{}
	if err := json.Unmarshal([]byte(jsondata), &seo); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewseoService(NewseoRepo())
	seo1, err := service.Create(seo)
	res, err := service.Delete(seo1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
