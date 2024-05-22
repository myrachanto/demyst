package gift

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	GiftService GiftServiceInterface = &giftService{}
)

type GiftServiceInterface interface {
	Create(Gift *Gift) (*Gift, httperrors.HttpErr)
	GetOne(code string) (Gift *Gift, errors httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, Gift *Gift) (string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type giftService struct {
	repo GiftrepoInterface
}

func NewGiftService(repository GiftrepoInterface) GiftServiceInterface {
	return &giftService{
		repository,
	}
}
func (service *giftService) Create(Gift *Gift) (*Gift, httperrors.HttpErr) {

	return service.repo.Create(Gift)

}

func (service *giftService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *giftService) GetOne(code string) (*Gift, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *giftService) Update(code string, Gift *Gift) (string, httperrors.HttpErr) {
	return service.repo.Update(code, Gift)
}
func (service *giftService) Delete(code string) (string, httperrors.HttpErr) {
	return service.repo.Delete(code)
}
