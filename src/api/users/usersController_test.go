package users

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// end to end testing
var (
	jsondatas = `{"fname":"Jane","lname":"Doe","uname":"doe","phone":"1234567","email":"email1@example.com","password":"1234567","address":"psd 456 king view"
	}`
	// jsondata1 = `{"email":   "email1@example.com","password": "1234567"}`
)

//make end to end testing
func TestCreateUser(t *testing.T) {
	user := &User{}
	if err := json.Unmarshal([]byte(jsondatas), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	// Bizname = "alias4567"
	u, err := service.Create(user)
	assert.EqualValues(t, "Jane", u.Firstname, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Usercode)
}

// func TestLoginUser(t *testing.T) {

// 	user1 := &User{}
// 	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
// 		t.Errorf("failed to unmarshal user data %v", err.Error())
// 	}
// 	_, _ = UserService.Create(user1)
// 	user := &LoginUser{}
// 	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
// 		t.Errorf("failed to unmarshal user data %v", err.Error())
// 	}
// 	res, err := UserService.Login(user)
// 	// expected := "user created successifully"
// 	if res.Token == "" {
// 		assert.EqualValues(t, "", res, "Something went wrong testing with the Logging method")
// 		assert.NotNil(t, err, "Something went wrong testing with the Logging method")
// 	}
// 	// afterparty cleaner
// 	//use uname to clean after testing
// 	UserService.Delete(res.Usercode)
// }
func TestGetAllUser(t *testing.T) {

	user1 := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	_, _ = service.Create(user1)
	_, err := service.GetAll("")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	// afterparty cleaner
	//use uname to clean after testing
	service.Delete(user1.Usercode)
}

func TestGetOneUser(t *testing.T) {

	user1 := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	person, _ := service.Create(user1)

	_, e := service.GetOne(person.Usercode)
	// assert.EqualValues(t, "doe", u.UName, "Something went wrong testing with the Getting one method")
	assert.Nil(t, e, "Something went wrong testing with the Getting one method")
	// afterparty cleaner
	//use uname to clean after testing
	service.Delete(user1.Usercode)
}

// func TestUpdateUser(t *testing.T) {
// 	user1 := &User{}
// 	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
// 		t.Errorf("failed to unmarshal user data %v", err.Error())
// 	}
// 	person, _ := UserService.Create(user1)
// 	// fmt.Println(">>>>>>>>>>>>>>sd", person)
// 	person.FName = "John"
// 	// fmt.Println(">>>>>>>>>>>>>>", person)
// 	u, err := UserService.Update(person.Usercode, person)
// 	// expected := "user created successifully"
// 	assert.EqualValues(t, "John", u.FName, "Something went wrong testing with the Getting one method")
// 	assert.Nil(t, err)
// 	// afterparty cleaner
// 	//use uname to clean after testing
// 	UserService.Delete(user1.Usercode)
// }
func TestDeleteUser(t *testing.T) {

	user1 := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewUserService(NewUserRepo())
	person, _ := service.Create(user1)
	res, err := service.Delete(person.Usercode)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
