package newssections

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	NewsSectionService NewsSectionServiceInterface = &newssectionsService{}
)

type NewsSectionServiceInterface interface {
	Update(code string, newssections *NewsSection) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type newssectionsService struct {
	repo NewsSectionrepoInterface
}

func NewnewssectionsService(repository NewsSectionrepoInterface) NewsSectionServiceInterface {
	return &newssectionsService{
		repository,
	}
}
func (service *newssectionsService) Update(code string, newssections *NewsSection) (string, httperrors.HttpErr) {
	return service.repo.Update(code, newssections)
}
func (service *newssectionsService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
