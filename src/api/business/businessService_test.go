package business

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"wed3455","year_established":2020}`

func TestServiceCreateBusiness(t *testing.T) {
	business := &Business{}
	if err := json.Unmarshal([]byte(jsondata), &business); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewbusinessService(NewbusinessRepo())
	u, err := service.Create(business)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllBusiness(t *testing.T) {

	business := &Business{}
	if err := json.Unmarshal([]byte(jsondata), &business); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewbusinessService(NewbusinessRepo())
	business1, _ := service.Create(business)
	_, err := service.GetAll("")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(business1.Code)
}
func TestServiceGetOneBusiness(t *testing.T) {

	business := &Business{}
	if err := json.Unmarshal([]byte(jsondata), &business); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewbusinessService(NewbusinessRepo())
	business1, err := service.Create(business)
	assert.EqualValues(t, "name", business.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(business1.Code)
}

func TestServiceUpdateBusiness(t *testing.T) {

	business := &Business{}
	if err := json.Unmarshal([]byte(jsondata), &business); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewbusinessService(NewbusinessRepo())
	business1, err := service.Create(business)
	assert.Nil(t, err)
	business1.Name = "updatedname"
	res, err := service.Update(business1.Code, business1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(business1.Code)
}
func TestServiceDeleteBusiness(t *testing.T) {

	business := &Business{}
	if err := json.Unmarshal([]byte(jsondata), &business); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewbusinessService(NewbusinessRepo())
	business1, err := service.Create(business)
	res, err := service.Delete(business1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
