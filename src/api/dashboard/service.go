package dashboard

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	DashboardService DashboardServiceInterface = &dashboardService{}
)

type DashboardServiceInterface interface {
	HomeCms() (*Dashboard, httperrors.HttpErr)
	Index() (*Home, httperrors.HttpErr)
	Layout() (*Nav, httperrors.HttpErr)
}
type dashboardService struct {
	repo DashboardrepoInterface
}

func NewdashboardService(repository DashboardrepoInterface) DashboardServiceInterface {
	return &dashboardService{
		repository,
	}
}
func (service *dashboardService) HomeCms() (*Dashboard, httperrors.HttpErr) {

	return service.repo.HomeCms()

}
func (service *dashboardService) Index() (*Home, httperrors.HttpErr) {

	return service.repo.Index()

}
func (service *dashboardService) Layout() (*Nav, httperrors.HttpErr) {

	return service.repo.Layout()

}
