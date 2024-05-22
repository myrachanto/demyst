package tags

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	TagService TagServiceInterface = &tagService{}
)

type TagServiceInterface interface {
	Create(tag *Tag) (*Tag, httperrors.HttpErr)
	GetOne(code string) (tag *Tag, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, tag *Tag) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type tagService struct {
	repo TagrepoInterface
}

func NewtagService(repository TagrepoInterface) TagServiceInterface {
	return &tagService{
		repository,
	}
}
func (service *tagService) Create(tag *Tag) (*Tag, httperrors.HttpErr) {

	return service.repo.Create(tag)

}

func (service *tagService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}

func (service *tagService) GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr) {
	tags, err := service.repo.GetAll1(searcher)
	return tags, err
}
func (service *tagService) GetOne(code string) (*Tag, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *tagService) Update(code string, tag *Tag) (string, httperrors.HttpErr) {
	return service.repo.Update(code, tag)
}
func (service *tagService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
