package rating

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	RatingService RatingServiceInterface = &ratingService{}
)

type RatingServiceInterface interface {
	Create(seo *Rating) (*Rating, httperrors.HttpErr)
	GetOne(id string) (*Rating, httperrors.HttpErr)
	GetAll() ([]*Rating, httperrors.HttpErr)
	Featured(code string, status bool) httperrors.HttpErr
	Update(id string, rating *Rating) (*Rating, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
}
type ratingService struct {
	repo RatingRepoInterface
}

func NewratingService(repository RatingRepoInterface) RatingServiceInterface {
	return &ratingService{
		repository,
	}
}

func (service *ratingService) Create(rating *Rating) (*Rating, httperrors.HttpErr) {
	rating, err1 := service.repo.Create(rating)
	return rating, err1

}

func (service *ratingService) Featured(code string, status bool) httperrors.HttpErr {
	err1 := service.repo.Featured(code, status)
	return err1
}
func (service *ratingService) GetOne(id string) (*Rating, httperrors.HttpErr) {
	rating, err1 := service.repo.GetOne(id)
	return rating, err1
}

func (service *ratingService) GetAll() ([]*Rating, httperrors.HttpErr) {
	ratings, err := service.repo.GetAll()
	return ratings, err
}

func (service *ratingService) Update(id string, rating *Rating) (*Rating, httperrors.HttpErr) {
	rating, err1 := service.repo.Update(id, rating)
	return rating, err1
}
func (service *ratingService) Delete(id string) (string, httperrors.HttpErr) {
	success, failure := service.repo.Delete(id)
	return success, failure
}
