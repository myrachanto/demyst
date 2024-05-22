package majorcategory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// end to end testing
var (
	jsondata = `{"name":"name","title":"title","description":"description","shopalias":"shopalias"}`
)

// make end to end testing
func TestCreatemajorcategory(t *testing.T) {
	majorcategory := &Majorcategory{}
	if err := json.Unmarshal([]byte(jsondata), &majorcategory); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	u, err := MajorcategoryService.Create(majorcategory)
	// fmt.Println(">>>>>>>>>", u)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	MajorcategoryService.Delete(u.Code)
}

func TestGetAllmajorcategory(t *testing.T) {

	majorcategory := &Majorcategory{}
	if err := json.Unmarshal([]byte(jsondata), &majorcategory); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	_, _ = MajorcategoryService.Create(majorcategory)
	_, err := MajorcategoryService.GetAll()
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	MajorcategoryService.Delete(majorcategory.Code)
}
func TestGetOnemajorcategory(t *testing.T) {

	majorcategory := &Majorcategory{}
	if err := json.Unmarshal([]byte(jsondata), &majorcategory); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	majorcategory, err := MajorcategoryService.Create(majorcategory)
	assert.EqualValues(t, "name", majorcategory.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	MajorcategoryService.Delete(majorcategory.Code)
}

func TestUpdatemajorcategory(t *testing.T) {
	majorcategory := &Majorcategory{}
	if err := json.Unmarshal([]byte(jsondata), &majorcategory); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	majorcategory1, _ := MajorcategoryService.Create(majorcategory)
	majorcategory1.Name = "update"
	u, err := MajorcategoryService.Update(majorcategory1.Code, majorcategory1)
	assert.EqualValues(t, "update", u.Name, "Something went wrong testing with the Getting one method")
	assert.Nil(t, err)
	MajorcategoryService.Delete(majorcategory1.Code)
}
func TestDeletemajorcategory(t *testing.T) {

	majorcategory1 := &Majorcategory{}
	if err := json.Unmarshal([]byte(jsondata), &majorcategory1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	majorcategory, _ := MajorcategoryService.Create(majorcategory1)
	res, err := MajorcategoryService.Delete(majorcategory.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
