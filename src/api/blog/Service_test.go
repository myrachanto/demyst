package blog

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/estate/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreateBlog(t *testing.T) {
	Blog := &Blog{}
	if err := json.Unmarshal([]byte(jsondata), &Blog); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewBlogService(NewBlogRepo())
	u, err := service.Create(Blog)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllBlog(t *testing.T) {

	Blog := &Blog{}
	if err := json.Unmarshal([]byte(jsondata), &Blog); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewBlogService(NewBlogRepo())
	Blog1, _ := service.Create(Blog)
	_, err := service.GetAll(support.Paginator{})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(Blog1.Code)
}
func TestServiceGetOneBlog(t *testing.T) {

	Blog := &Blog{}
	if err := json.Unmarshal([]byte(jsondata), &Blog); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewBlogService(NewBlogRepo())
	Blog1, err := service.Create(Blog)
	assert.EqualValues(t, "name", Blog.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(Blog1.Code)
}

func TestServiceUpdateBlog(t *testing.T) {

	Blog := &Blog{}
	if err := json.Unmarshal([]byte(jsondata), &Blog); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewBlogService(NewBlogRepo())
	Blog1, err := service.Create(Blog)
	assert.Nil(t, err)
	Blog1.Name = "updatedname"
	res, err := service.Update(Blog1.Code, Blog1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(Blog1.Code)
}
func TestServiceDeleteBlog(t *testing.T) {

	Blog := &Blog{}
	if err := json.Unmarshal([]byte(jsondata), &Blog); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewBlogService(NewBlogRepo())
	Blog1, err := service.Create(Blog)
	res, err := service.Delete(Blog1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
