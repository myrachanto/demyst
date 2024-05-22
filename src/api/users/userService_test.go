package users

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/estate/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata1 = `{"fullname":"jane","username":"doe","Usercode": "Doe345","Phone":"1234567","Email":   "email@example.com","Password": "1234567","Address":"psd 456 king view"}`

func TestServiceCreateUser(t *testing.T) {
	user := &User{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	u, err := service.Create(user)
	// assert.EqualValues(t, "jane", u.Fullname)
	assert.Nil(t, err)
	service.Delete(u.Usercode)
}
func TestServiceGetAllUser(t *testing.T) {

	user := &User{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	user1, _ := service.Create(user)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err)
	service.Delete(user1.Usercode)
}
func TestServiceGetOneUser(t *testing.T) {

	user := &User{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	user1, err := service.Create(user)
	assert.Nil(t, err)
	service.Delete(user1.Usercode)
}

func TestServiceDeleteUser(t *testing.T) {

	user := &User{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	user1, err := service.Create(user)
	assert.Nil(t, err)
	res, err := service.Delete(user1.Usercode)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res)
	assert.Nil(t, err)
}
