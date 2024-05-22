package profilesections

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	ProfileSectionService ProfileSectionServiceInterface = &ProfilesectionsService{}
)

type ProfileSectionServiceInterface interface {
	Update(code string, Profilesections *ProfileSection) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type ProfilesectionsService struct {
	repo ProfileSectionrepoInterface
}

func NewProfilesectionsService(repository ProfileSectionrepoInterface) ProfileSectionServiceInterface {
	return &ProfilesectionsService{
		repository,
	}
}
func (service *ProfilesectionsService) Update(code string, Profilesections *ProfileSection) (string, httperrors.HttpErr) {
	return service.repo.Update(code, Profilesections)
}
func (service *ProfilesectionsService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
