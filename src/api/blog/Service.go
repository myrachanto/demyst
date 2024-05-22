package blog

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	BlogService BlogServiceInterface = &blogService{}
)

type BlogServiceInterface interface {
	Create(Blog *Blog) (*Blog, httperrors.HttpErr)
	GetOne(code string) (Blog *Blog, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, Blog *Blog) (string, httperrors.HttpErr)
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateTrending(code string, status bool) httperrors.HttpErr
	UpdateExclusive(code string, status bool) httperrors.HttpErr
	Delete(code string) (string, httperrors.HttpErr)
	GetOneByUrl(code string) (*ByBlog, httperrors.HttpErr)
	GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr)
	CreateComment(coment *Comment) httperrors.HttpErr
	DeleteComment(code, Blogcode string) httperrors.HttpErr
}
type blogService struct {
	repo BlogRepoInterface
}

func NewBlogService(repository BlogRepoInterface) BlogServiceInterface {
	return &blogService{
		repository,
	}
}
func (service *blogService) Create(Blog *Blog) (*Blog, httperrors.HttpErr) {
	return service.repo.Create(Blog)
}
func (service *blogService) CreateComment(coment *Comment) httperrors.HttpErr {
	return service.repo.CreateComment(coment)
}
func (service *blogService) DeleteComment(code, Blogcode string) httperrors.HttpErr {
	return service.repo.DeleteComment(code, Blogcode)
}

func (service *blogService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *blogService) GetOne(code string) (*Blog, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *blogService) GetOneByUrl(code string) (*ByBlog, httperrors.HttpErr) {
	return service.repo.GetOneByUrl(code)
}
func (service *blogService) GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetByCategory(code, search)
}
func (service *blogService) Update(code string, Blog *Blog) (string, httperrors.HttpErr) {
	return service.repo.Update(code, Blog)
}
func (service *blogService) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateFeatured(code, status)
}
func (service *blogService) UpdateTrending(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateTrending(code, status)
}
func (service *blogService) UpdateExclusive(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateExclusive(code, status)
}
func (service *blogService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
