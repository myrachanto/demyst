package business

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// businessController ...
var (
	BusinessController BusinessControllerInterface = &businessController{}
)

type BusinessControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type businessController struct {
	service businessServiceInterface
}

func NewbusinessController(ser businessServiceInterface) BusinessControllerInterface {
	return &businessController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a business
// @Description Create a new business item
// @Tags businesss
// @Accept json
// @Produce json
// @Success 201 {object} business
// @Failure 400 {object} support.HttpError
// @Router /api/businesss [post]
func (controller businessController) Create(c echo.Context) error {

	business := &Business{}
	fmt.Println("--------------------create")
	business.Name = c.FormValue("name")
	business.BusinessPin = c.FormValue("business_pin")
	yof, err := strconv.ParseUint(c.FormValue("year_established"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusCreated, "failed to parse the year of establishment")
	}
	business.YearEstablished = int(yof)
	_, err1 := controller.service.Create(business)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "business created succesifully")

}

// GetAll godoc
// @Summary GetAll a business
// @Description Getall businesss
// @Tags businesss
// @Accept json
// @Produce json
// @Success 201 {object} business
// @Failure 400 {object} support.HttpError
// @Router /api/businesss [get]
func (controller businessController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")
	businesss, err3 := controller.service.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, businesss)
}

// @Summary Get a business
// @Description Get item
// @Tags businesss
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} business
// @Failure 400 {object} support.HttpError
// @Router /api/businesss [get]
func (controller businessController) GetOne(c echo.Context) error {
	code := c.Param("code")
	business, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, business)
}

// Create godoc
// @Summary Update an business
// @Description Update an new business item
// @Tags businesss
// @Accept json
// @Produce json
// @Success 200 {object} business
// @Failure 400 {object} support.HttpError
// @Router /api/businesss [put]
func (controller businessController) Update(c echo.Context) error {

	business := &Business{}
	business.Name = c.FormValue("name")
	business.BusinessPin = c.FormValue("business_pin")
	yof, err := strconv.ParseUint(c.FormValue("year_established"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusCreated, "failed to parse the loan amount")
	}
	business.YearEstablished = int(yof)
	code := c.Param("code")
	_, err1 := controller.service.Update(code, business)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "business created succesifully")
}

// @Summary Delte a business
// @Description Delte item
// @Tags businesss
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} business
// @Failure 400 {object} support.HttpError
// @Router /api/businesss [delete]
func (controller businessController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
