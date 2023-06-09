package pages

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/support"
)

var (
	PageService PageServiceInterface = &pageService{}
)

type PageServiceInterface interface {
	Create(page *Page) (*Page, httperrors.HttpErr)
	GetOne(code string) (page *Page, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, page *Page) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
	GetOneByUrl(code string) (*Page, httperrors.HttpErr)
}
type pageService struct {
	repo PagerepoInterface
}

func NewpageService(repository PagerepoInterface) PageServiceInterface {
	return &pageService{
		repository,
	}
}
func (service *pageService) Create(page *Page) (*Page, httperrors.HttpErr) {

	return service.repo.Create(page)

}

func (service *pageService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *pageService) GetOne(code string) (*Page, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *pageService) GetOneByUrl(code string) (*Page, httperrors.HttpErr) {
	return service.repo.GetOneByUrl(code)
}
func (service *pageService) Update(code string, page *Page) (string, httperrors.HttpErr) {
	return service.repo.Update(code, page)
}
func (service *pageService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
