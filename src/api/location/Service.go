package location

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	LocationService LocationServiceInterface = &locationService{}
)

type LocationServiceInterface interface {
	Create(location *Location) (*Location, httperrors.HttpErr)
	GetOne(code string) (location *Location, errors httperrors.HttpErr)
	GetAll() ([]Location, httperrors.HttpErr)
	GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, location *Location) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type locationService struct {
	repo locationrepoInterface
}

func NewlocationService(repository locationrepoInterface) LocationServiceInterface {
	return &locationService{
		repository,
	}
}
func (service *locationService) Create(location *Location) (*Location, httperrors.HttpErr) {

	return service.repo.Create(location)

}

func (service *locationService) GetAll() ([]Location, httperrors.HttpErr) {
	return service.repo.GetAll()
}

func (service *locationService) GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr) {
	tags, err := service.repo.GetAll1(searcher)
	return tags, err
}
func (service *locationService) GetOne(code string) (*Location, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *locationService) Update(code string, location *Location) (string, httperrors.HttpErr) {
	return service.repo.Update(code, location)
}
func (service *locationService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
