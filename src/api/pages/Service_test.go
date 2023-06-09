package pages

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/sports/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatepage(t *testing.T) {
	page := &Page{}
	if err := json.Unmarshal([]byte(jsondata), &page); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewpageService(NewpageRepo())
	u, err := service.Create(page)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllpage(t *testing.T) {

	page := &Page{}
	if err := json.Unmarshal([]byte(jsondata), &page); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewpageService(NewpageRepo())
	page1, _ := service.Create(page)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(page1.Code)
}
func TestServiceGetOnepage(t *testing.T) {

	page := &Page{}
	if err := json.Unmarshal([]byte(jsondata), &page); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewpageService(NewpageRepo())
	page1, err := service.Create(page)
	assert.EqualValues(t, "name", page.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(page1.Code)
}

func TestServiceUpdatepage(t *testing.T) {

	page := &Page{}
	if err := json.Unmarshal([]byte(jsondata), &page); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewpageService(NewpageRepo())
	page1, err := service.Create(page)
	assert.Nil(t, err)
	page1.Name = "updatedname"
	res, err := service.Update(page1.Code, page1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(page1.Code)
}
func TestServiceDeletepage(t *testing.T) {

	page := &Page{}
	if err := json.Unmarshal([]byte(jsondata), &page); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewpageService(NewpageRepo())
	page1, err := service.Create(page)
	res, err := service.Delete(page1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
