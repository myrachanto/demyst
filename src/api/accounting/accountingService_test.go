package accounting

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreateAccounting(t *testing.T) {
	accounting := &Accounting{}
	if err := json.Unmarshal([]byte(jsondata), &accounting); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewaccountingService(NewaccountingRepo())
	u, err := service.Create(accounting)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllAccounting(t *testing.T) {

	accounting := &Accounting{}
	if err := json.Unmarshal([]byte(jsondata), &accounting); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewaccountingService(NewaccountingRepo())
	accounting1, _ := service.Create(accounting)
	_, err := service.GetAll("")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(accounting1.Code)
}
func TestServiceGetOneAccounting(t *testing.T) {

	accounting := &Accounting{}
	if err := json.Unmarshal([]byte(jsondata), &accounting); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewaccountingService(NewaccountingRepo())
	accounting1, err := service.Create(accounting)
	assert.EqualValues(t, "name", accounting.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(accounting1.Code)
}

func TestServiceUpdateAccounting(t *testing.T) {

	accounting := &Accounting{}
	if err := json.Unmarshal([]byte(jsondata), &accounting); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewaccountingService(NewaccountingRepo())
	accounting1, err := service.Create(accounting)
	assert.Nil(t, err)
	accounting1.Name = "updatedname"
	res, err := service.Update(accounting1.Code, accounting1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(accounting1.Code)
}
func TestServiceDeleteAccounting(t *testing.T) {

	accounting := &Accounting{}
	if err := json.Unmarshal([]byte(jsondata), &accounting); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewaccountingService(NewaccountingRepo())
	accounting1, err := service.Create(accounting)
	res, err := service.Delete(accounting1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
