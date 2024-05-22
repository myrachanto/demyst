package profile

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	ProfileService ProfileServiceInterface = &profileService{}
)

type ProfileServiceInterface interface {
	Create(Profile *Profile) (*Profile, httperrors.HttpErr)
	GetOne(code string) (Profile *Profile, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, Profile *Profile) (string, httperrors.HttpErr)
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateTrending(code string, status bool) httperrors.HttpErr
	UpdateExclusive(code string, status bool) httperrors.HttpErr
	Delete(code string) (string, httperrors.HttpErr)
	GetOneByUrl(code string) (*ByProfile, httperrors.HttpErr)
	GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr)
	CreateComment(coment *Comment) httperrors.HttpErr
	DeleteComment(code, Profilecode string) httperrors.HttpErr
}
type profileService struct {
	repo ProfileRepoInterface
}

func NewprofileService(repository ProfileRepoInterface) ProfileServiceInterface {
	return &profileService{
		repository,
	}
}
func (service *profileService) Create(Profile *Profile) (*Profile, httperrors.HttpErr) {
	return service.repo.Create(Profile)
}
func (service *profileService) CreateComment(coment *Comment) httperrors.HttpErr {
	return service.repo.CreateComment(coment)
}
func (service *profileService) DeleteComment(code, Profilecode string) httperrors.HttpErr {
	return service.repo.DeleteComment(code, Profilecode)
}

func (service *profileService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *profileService) GetOne(code string) (*Profile, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *profileService) GetOneByUrl(code string) (*ByProfile, httperrors.HttpErr) {
	return service.repo.GetOneByUrl(code)
}
func (service *profileService) GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetByCategory(code, search)
}
func (service *profileService) Update(code string, Profile *Profile) (string, httperrors.HttpErr) {
	return service.repo.Update(code, Profile)
}
func (service *profileService) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateFeatured(code, status)
}
func (service *profileService) UpdateTrending(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateTrending(code, status)
}
func (service *profileService) UpdateExclusive(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateExclusive(code, status)
}
func (service *profileService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
