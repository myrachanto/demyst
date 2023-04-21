package accounting

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	AccountingService AccountingServiceInterface = &accountingService{}
)

type AccountingServiceInterface interface {
	Create(accounting *Accounting) (*Accounting, httperrors.HttpErr)
	GetOne(code string) (accounting *Accounting, errors httperrors.HttpErr)
	GetAll(search string) ([]*Accounting, httperrors.HttpErr)
	Update(code string, accounting *Accounting) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type accountingService struct {
	repo AccountingrepoInterface
}

func NewaccountingService(repository AccountingrepoInterface) AccountingServiceInterface {
	return &accountingService{
		repository,
	}
}
func (service *accountingService) Create(accounting *Accounting) (*Accounting, httperrors.HttpErr) {

	return service.repo.Create(accounting)

}

func (service *accountingService) GetAll(search string) ([]*Accounting, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *accountingService) GetOne(code string) (*Accounting, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *accountingService) Update(code string, accounting *Accounting) (string, httperrors.HttpErr) {
	return service.repo.Update(code, accounting)
}
func (service *accountingService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
