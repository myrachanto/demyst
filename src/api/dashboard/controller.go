package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// categoryController ...
var (
	DashboardController DashboardControllerInterface = &dashboardController{}
)

type DashboardControllerInterface interface {
	HomeCms(c echo.Context) error
	Index(c echo.Context) error
}

type dashboardController struct {
	service DashboardServiceInterface
}

func NewdashboardController(ser DashboardServiceInterface) DashboardControllerInterface {
	return &dashboardController{
		ser,
	}
}

// @Summary  Get  Dashboard
// @Description  Get  Dashboard
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 201 {object} Dashboard
// @Failure 400 {object} support.HttpError
// @Router /api/dashboard [get]
func (controller dashboardController) HomeCms(c echo.Context) error {
	category, problem := controller.service.HomeCms()
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, category)
}

// @Summary  Get  Home
// @Description  Get  Home
// @Tags Home
// @Accept json
// @Produce json
// @Success 201 {object} Home
// @Failure 400 {object} support.HttpError
// @Router /home [get]
func (controller dashboardController) Index(c echo.Context) error {
	category, problem := controller.service.Index()
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, category)
}
