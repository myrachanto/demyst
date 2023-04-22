package business

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	BusinessService businessServiceInterface = &businessService{}
)

type businessServiceInterface interface {
	Create(business *Business) (*Business, httperrors.HttpErr)
	GetOne(code string) (business *Business, errors httperrors.HttpErr)
	GetAll(search string) ([]*Business, httperrors.HttpErr)
	Update(code string, business *Business) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type businessService struct {
	repo BusinessrepoInterface
}

func NewbusinessService(repository BusinessrepoInterface) businessServiceInterface {
	return &businessService{
		repository,
	}
}
func (service *businessService) Create(business *Business) (*Business, httperrors.HttpErr) {

	return service.repo.Create(business)

}

func (service *businessService) GetAll(search string) ([]*Business, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *businessService) GetOne(code string) (*Business, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *businessService) Update(code string, business *Business) (string, httperrors.HttpErr) {
	return service.repo.Update(code, business)
}
func (service *businessService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
