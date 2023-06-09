package category

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/support"
)

var (
	CategoryService CategoryServiceInterface = &categoryService{}
)

type CategoryServiceInterface interface {
	Create(category *Category) (*Category, httperrors.HttpErr)
	GetOne(code string) (category *Category, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, category *Category) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type categoryService struct {
	repo CategoryrepoInterface
}

func NewcategoryService(repository CategoryrepoInterface) CategoryServiceInterface {
	return &categoryService{
		repository,
	}
}
func (service *categoryService) Create(category *Category) (*Category, httperrors.HttpErr) {

	return service.repo.Create(category)

}

func (service *categoryService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *categoryService) GetOne(code string) (*Category, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *categoryService) Update(code string, category *Category) (string, httperrors.HttpErr) {
	return service.repo.Update(code, category)
}
func (service *categoryService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
