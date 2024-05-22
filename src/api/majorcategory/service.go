package majorcategory

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	MajorcategoryService MajorcategoryServiceInterface = &majorcategoryService{}
)

type MajorcategoryServiceInterface interface {
	Create(seo *Majorcategory) (*Majorcategory, httperrors.HttpErr)
	GetOne(id string) (*Majorcategory, httperrors.HttpErr)
	GetAll() ([]Majorcategory, httperrors.HttpErr)
	GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr)
	Update(id string, majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
}
type majorcategoryService struct {
	repo MajorcategoryRepoInterface
}

func NewmajorcategoryService(repository MajorcategoryRepoInterface) MajorcategoryServiceInterface {
	return &majorcategoryService{
		repository,
	}
}

func (service *majorcategoryService) Create(majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr) {
	res, err1 := service.repo.Create(majorcategory)
	return res, err1

}

func (service *majorcategoryService) GetOne(id string) (*Majorcategory, httperrors.HttpErr) {
	majorcategory, err1 := service.repo.GetOne(id)
	return majorcategory, err1
}

func (service *majorcategoryService) GetAll() ([]Majorcategory, httperrors.HttpErr) {
	majorcategorys, err := service.repo.GetAll()
	return majorcategorys, err
}
func (service *majorcategoryService) GetAll1(searcher support.Paginator) (*Results, httperrors.HttpErr) {
	tags, err := service.repo.GetAll1(searcher)
	return tags, err
}

func (service *majorcategoryService) Update(id string, majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr) {
	res, err1 := service.repo.Update(id, majorcategory)
	return res, err1
}
func (service *majorcategoryService) Delete(id string) (string, httperrors.HttpErr) {
	success, failure := service.repo.Delete(id)
	return success, failure
}
