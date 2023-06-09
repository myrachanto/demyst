package news

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/sports/src/support"
	"github.com/stretchr/testify/assert"
)

var jsondata = `{"name":"name","business_pin":"afsd3455","urlEndpoint":"https://example.com"}`

func TestServiceCreatenews(t *testing.T) {
	news := &News{}
	if err := json.Unmarshal([]byte(jsondata), &news); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewnewsService(NewnewsRepo())
	u, err := service.Create(news)
	assert.EqualValues(t, "name", u.Name, "failed to validate create method")
	assert.Nil(t, err)
	service.Delete(u.Code)
}
func TestServiceGetAllnews(t *testing.T) {

	news := &News{}
	if err := json.Unmarshal([]byte(jsondata), &news); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewnewsService(NewnewsRepo())
	news1, _ := service.Create(news)
	_, err := service.GetAll(support.Paginator{"", "", 0, 0})
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	service.Delete(news1.Code)
}
func TestServiceGetOnenews(t *testing.T) {

	news := &News{}
	if err := json.Unmarshal([]byte(jsondata), &news); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewnewsService(NewnewsRepo())
	news1, err := service.Create(news)
	assert.EqualValues(t, "name", news.Name, "Something went wrong testing with the Getting all method")
	assert.Nil(t, err, "Something went wrong testing with the Getting one method")
	service.Delete(news1.Code)
}

func TestServiceUpdatenews(t *testing.T) {

	news := &News{}
	if err := json.Unmarshal([]byte(jsondata), &news); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewnewsService(NewnewsRepo())
	news1, err := service.Create(news)
	assert.Nil(t, err)
	news1.Name = "updatedname"
	res, err := service.Update(news1.Code, news1)
	assert.EqualValues(t, "successifully Updated!", res, "Something went wrong testing with the updatting")
	assert.Nil(t, err)
	service.Delete(news1.Code)
}
func TestServiceDeletenews(t *testing.T) {

	news := &News{}
	if err := json.Unmarshal([]byte(jsondata), &news); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	service := NewnewsService(NewnewsRepo())
	news1, err := service.Create(news)
	res, err := service.Delete(news1.Code)
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
