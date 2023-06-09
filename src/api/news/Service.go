package news

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/support"
)

var (
	NewsService NewsServiceInterface = &newsService{}
)

type NewsServiceInterface interface {
	Create(news *News) (*News, httperrors.HttpErr)
	GetOne(code string) (news *News, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, news *News) (string, httperrors.HttpErr)
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateTrending(code string, status bool) httperrors.HttpErr
	UpdateExclusive(code string, status bool) httperrors.HttpErr
	Delete(code string) (string, httperrors.HttpErr)
	GetOneByUrl(code string) (*ByNews, httperrors.HttpErr)
	GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr)
}
type newsService struct {
	repo NewsrepoInterface
}

func NewnewsService(repository NewsrepoInterface) NewsServiceInterface {
	return &newsService{
		repository,
	}
}
func (service *newsService) Create(news *News) (*News, httperrors.HttpErr) {

	return service.repo.Create(news)

}

func (service *newsService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *newsService) GetOne(code string) (*News, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *newsService) GetOneByUrl(code string) (*ByNews, httperrors.HttpErr) {
	return service.repo.GetOneByUrl(code)
}
func (service *newsService) GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetByCategory(code, search)
}
func (service *newsService) Update(code string, news *News) (string, httperrors.HttpErr) {
	return service.repo.Update(code, news)
}
func (service *newsService) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateFeatured(code, status)
}
func (service *newsService) UpdateTrending(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateTrending(code, status)
}
func (service *newsService) UpdateExclusive(code string, status bool) httperrors.HttpErr {
	return service.repo.UpdateExclusive(code, status)
}
func (service *newsService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
