package seo

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	SeoService SeoServiceInterface = &seoService{}
)

type SeoServiceInterface interface {
	Create(seo *Seo) (*Seo, httperrors.HttpErr)
	GetOne(code string) (seo *Seo, errors httperrors.HttpErr)
	GetAll(search string) ([]Seo, httperrors.HttpErr)
	GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, seo *Seo) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type seoService struct {
	repo SeorepoInterface
}

func NewseoService(repository SeorepoInterface) SeorepoInterface {
	return &seoService{
		repository,
	}
}
func (service *seoService) Create(seo *Seo) (*Seo, httperrors.HttpErr) {

	return service.repo.Create(seo)

}

func (service *seoService) GetAll(search string) ([]Seo, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}

func (service *seoService) GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr) {
	tags, err := service.repo.GetAll1(searcher)
	return tags, err
}
func (service *seoService) GetOne(code string) (*Seo, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *seoService) Update(code string, seo *Seo) (string, httperrors.HttpErr) {
	return service.repo.Update(code, seo)
}
func (service *seoService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
