package subLocation

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	SubLocationService SubLocationServiceInterface = &subLocationService{}
)

type SubLocationServiceInterface interface {
	Create(subLocation *SubLocation) (*SubLocation, httperrors.HttpErr)
	GetOne(code string) (subLocation *SubLocation, errors httperrors.HttpErr)
	GetAll() ([]SubLocation, httperrors.HttpErr)
	GetAllByLocation(code string) ([]SubLocation, httperrors.HttpErr)
	GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, subLocation *SubLocation) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type subLocationService struct {
	repo SubLocationrepoInterface
}

func NewsubLocationService(repository SubLocationrepoInterface) SubLocationServiceInterface {
	return &subLocationService{
		repository,
	}
}
func (service *subLocationService) Create(subLocation *SubLocation) (*SubLocation, httperrors.HttpErr) {

	return service.repo.Create(subLocation)

}

func (service *subLocationService) GetAll() ([]SubLocation, httperrors.HttpErr) {
	return service.repo.GetAll()
}

func (service *subLocationService) GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr) {
	tags, err := service.repo.GetAll1(searcher)
	return tags, err
}
func (service *subLocationService) GetAllByLocation(code string) ([]SubLocation, httperrors.HttpErr) {
	tags, err := service.repo.GetAllByLocation(code)
	return tags, err
}
func (service *subLocationService) GetOne(code string) (*SubLocation, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *subLocationService) Update(code string, subLocation *SubLocation) (string, httperrors.HttpErr) {
	return service.repo.Update(code, subLocation)
}
func (service *subLocationService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
