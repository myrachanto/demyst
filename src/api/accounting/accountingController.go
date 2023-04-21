package accounting

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// accountingController ...
var (
	AccountingController AccountingControllerInterface = &accountingController{}
)

type AccountingControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type accountingController struct {
	service AccountingServiceInterface
}

func NewaccountingController(ser AccountingServiceInterface) AccountingControllerInterface {
	return &accountingController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a accounting
// @Description Create a new accounting item
// @Tags accountings
// @Accept json
// @Produce json
// @Success 201 {object} accounting
// @Failure 400 {object} support.HttpError
// @Router /api/accountings [post]
func (controller accountingController) Create(c echo.Context) error {

	accounting := &Accounting{}
	fmt.Println("--------------------create")
	accounting.Name = c.FormValue("name")
	accounting.BusinessPin = c.FormValue("business_pin")
	accounting.UrlEndpoint = c.FormValue("urlEndpoint")
	_, err1 := controller.service.Create(accounting)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "accounting created succesifully")

}

// GetAll godoc
// @Summary GetAll a accounting
// @Description Getall accountings
// @Tags accountings
// @Accept json
// @Produce json
// @Success 201 {object} accounting
// @Failure 400 {object} support.HttpError
// @Router /api/accountings [get]
func (controller accountingController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")
	accountings, err3 := controller.service.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, accountings)
}

// @Summary Get a accounting
// @Description Get item
// @Tags accountings
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} accounting
// @Failure 400 {object} support.HttpError
// @Router /api/accountings [get]
func (controller accountingController) GetOne(c echo.Context) error {
	code := c.Param("code")
	accounting, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, accounting)
}

// Create godoc
// @Summary Update an accounting
// @Description Update an new accounting item
// @Tags accountings
// @Accept json
// @Produce json
// @Success 200 {object} accounting
// @Failure 400 {object} support.HttpError
// @Router /api/accountings [put]
func (controller accountingController) Update(c echo.Context) error {

	accounting := &Accounting{}
	accounting.Name = c.FormValue("name")
	accounting.BusinessPin = c.FormValue("business_pin")
	accounting.UrlEndpoint = c.FormValue("urlEndpoint")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, accounting)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "accounting created succesifully")
}

// @Summary Delte a accounting
// @Description Delte item
// @Tags accountings
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} accounting
// @Failure 400 {object} support.HttpError
// @Router /api/accountings [delete]
func (controller accountingController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
