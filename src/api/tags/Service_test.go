package tags

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/estate/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatetag(t *testing.T) {
	tag := &Tag{}
	if err := json.Unmarshal([]byte(jsondata), &tag); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewtagService(NewtagRepo())
	u, err := service.Create(tag)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAlltag(t *testing.T) {

	tag := &Tag{}
	if err := json.Unmarshal([]byte(jsondata), &tag); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewtagService(NewtagRepo())
	tag1, _ := service.Create(tag)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(tag1.Code)
}
func TestServiceGetOnetag(t *testing.T) {

	tag := &Tag{}
	if err := json.Unmarshal([]byte(jsondata), &tag); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewtagService(NewtagRepo())
	tag1, err := service.Create(tag)
	assert.EqualValues(t, "name", tag.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(tag1.Code)
}

func TestServiceUpdatetag(t *testing.T) {

	tag := &Tag{}
	if err := json.Unmarshal([]byte(jsondata), &tag); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewtagService(NewtagRepo())
	tag1, err := service.Create(tag)
	assert.Nil(t, err)
	tag1.Name = "updatedname"
	res, err := service.Update(tag1.Code, tag1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(tag1.Code)
}
func TestServiceDeletetag(t *testing.T) {

	tag := &Tag{}
	if err := json.Unmarshal([]byte(jsondata), &tag); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewtagService(NewtagRepo())
	tag1, err := service.Create(tag)
	res, err := service.Delete(tag1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
