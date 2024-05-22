package gift

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/estate/src/support"
)

// GiftController ...
var (
	GiftController GiftControllerInterface = &giftController{}
)

type GiftControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type giftController struct {
	service GiftServiceInterface
}

func NewGiftController(ser GiftServiceInterface) GiftControllerInterface {
	return &giftController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a Gift
// @Description Create a new Gift item
// @Tags Gifts
// @Accept json
// @Produce json
// @Success 201 {object} Gift
// @Failure 400 {object} support.HttpError
// @Router /api/Gifts [post]

func (controller giftController) Create(c echo.Context) error {

	Gift := &Gift{}
	fmt.Println("--------------------create")
	Gift.Name = c.FormValue("name")
	Gift.Title = c.FormValue("title")
	Gift.Description = c.FormValue("description")
	_, err1 := controller.service.Create(Gift)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "Gift created succesifully")

}

// GetAll godoc
// @Summary GetAll a Gift
// @Description Getall Gifts
// @Tags Gifts
// @Accept json
// @Produce json
// @Success 201 {object} Gift
// @Failure 400 {object} support.HttpError
// @Router /api/Gifts [get]
func (controller giftController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh", search, ps, pn)
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator{Page: page, Pagesize: pagesize, Search: search}
	Gifts, err3 := controller.service.GetAll(searcher)

	// fmt.Println("----------------------results", Gifts)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, Gifts)
}

// GetOne godoc
// @Summary Get a Gift
// @Description Get item
// @Tags Gifts
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Gift
// @Failure 400 {object} support.HttpError
// @Router /api/Gifts [get]
func (controller giftController) GetOne(c echo.Context) error {
	code := c.Param("code")
	Gift, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Gift)
}

// Create godoc
// @Summary Update an Gift
// @Description Update an new Gift item
// @Tags Gifts
// @Accept json
// @Produce json
// @Success 200 {object} Gift
// @Failure 400 {object} support.HttpError
// @Router /api/Gifts [put]
func (controller giftController) Update(c echo.Context) error {

	Gift := &Gift{}
	Gift.Name = c.FormValue("name")
	Gift.Title = c.FormValue("title")
	Gift.Description = c.FormValue("description")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, Gift)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "Gift created succesifully")
}

// Delete godoc
// @Summary Delte a Gift
// @Description Delte item
// @Tags Gifts
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Gift
// @Failure 400 {object} support.HttpError
// @Router /api/Gifts [delete]
func (controller giftController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
