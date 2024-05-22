package BlogSections

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	BlogSectionService BlogSectionServiceInterface = &blogsectionsService{}
)

type BlogSectionServiceInterface interface {
	Update(code string, blogsections *BlogSection) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type blogsectionsService struct {
	repo BlogSectionrepoInterface
}

func NewblogsectionsService(repository BlogSectionrepoInterface) BlogSectionServiceInterface {
	return &blogsectionsService{
		repository,
	}
}
func (service *blogsectionsService) Update(code string, blogsections *BlogSection) (string, httperrors.HttpErr) {
	return service.repo.Update(code, blogsections)
}
func (service *blogsectionsService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
