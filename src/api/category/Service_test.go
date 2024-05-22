package category

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/estate/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatecategory(t *testing.T) {
	category := &Category{}
	if err := json.Unmarshal([]byte(jsondata), &category); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewcategoryService(NewcategoryRepo())
	u, err := service.Create(category)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllcategory(t *testing.T) {

	category := &Category{}
	if err := json.Unmarshal([]byte(jsondata), &category); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewcategoryService(NewcategoryRepo())
	category1, _ := service.Create(category)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(category1.Code)
}
func TestServiceGetOnecategory(t *testing.T) {

	category := &Category{}
	if err := json.Unmarshal([]byte(jsondata), &category); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewcategoryService(NewcategoryRepo())
	category1, err := service.Create(category)
	assert.EqualValues(t, "name", category.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(category1.Code)
}

func TestServiceUpdatecategory(t *testing.T) {

	category := &Category{}
	if err := json.Unmarshal([]byte(jsondata), &category); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewcategoryService(NewcategoryRepo())
	category1, err := service.Create(category)
	assert.Nil(t, err)
	category1.Name = "updatedname"
	res, err := service.Update(category1.Code, category1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(category1.Code)
}
func TestServiceDeletecategory(t *testing.T) {

	category := &Category{}
	if err := json.Unmarshal([]byte(jsondata), &category); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewcategoryService(NewcategoryRepo())
	category1, err := service.Create(category)
	res, err := service.Delete(category1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
