package feature

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	FeatureService FeatureServiceInterface = &featureService{}
)

type FeatureServiceInterface interface {
	Create(feature *Feature) (*Feature, httperrors.HttpErr)
	GetOne(code string) (feature *Feature, errors httperrors.HttpErr)
	GetAll(search string) ([]Feature, httperrors.HttpErr)
	GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, feature *Feature) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type featureService struct {
	repo FeaturerepoInterface
}

func NewfeatureService(repository FeaturerepoInterface) FeaturerepoInterface {
	return &featureService{
		repository,
	}
}
func (service *featureService) Create(feature *Feature) (*Feature, httperrors.HttpErr) {

	return service.repo.Create(feature)

}

func (service *featureService) GetAll(search string) ([]Feature, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}

func (service *featureService) GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr) {
	tags, err := service.repo.GetAll1(searcher)
	return tags, err
}
func (service *featureService) GetOne(code string) (*Feature, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *featureService) Update(code string, feature *Feature) (string, httperrors.HttpErr) {
	return service.repo.Update(code, feature)
}
func (service *featureService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
