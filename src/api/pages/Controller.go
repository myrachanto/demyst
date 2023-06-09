package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/sports/src/support"
)

// pageController ...
var (
	PageController PageControllerInterface = &pageController{}
)

type PageControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetOneByUrl(c echo.Context) error
}

type pageController struct {
	service PageServiceInterface
}

func NewpageController(ser PageServiceInterface) PageControllerInterface {
	return &pageController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a page
// @Description Create a new page item
// @Tags pages
// @Accept json
// @Produce json
// @Success 201 {object} Page
// @Failure 400 {object} support.HttpError
// @Router /api/pages [post]
func (controller pageController) Create(c echo.Context) error {

	page := &Page{}
	fmt.Println("--------------------create")
	page.Name = c.FormValue("name")
	page.Title = c.FormValue("title")
	page.Meta = c.FormValue("meta")
	page.Content = c.FormValue("content")
	_, err1 := controller.service.Create(page)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "page created succesifully")

}

// GetAll godoc
// @Summary GetAll a page
// @Description Getall pages
// @Tags pages
// @Accept json
// @Produce json
// @Success 201 {object} Page
// @Failure 400 {object} support.HttpError
// @Router /api/pages [get]
func (controller pageController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh")
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
	pages, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, pages)
}

// @Summary Get a page
// @Description Get item
// @Tags pages
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Page
// @Failure 400 {object} support.HttpError
// @Router /api/pages [get]
func (controller pageController) GetOne(c echo.Context) error {
	code := c.Param("code")
	page, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, page)
}

// @Summary Get a page
// @Description Get item
// @Tags pages
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Page
// @Failure 400 {object} support.HttpError
// @Router /api/pages [get]
func (controller pageController) GetOneByUrl(c echo.Context) error {
	code := c.Param("code")
	fmt.Println("================", code)
	page, problem := controller.service.GetOneByUrl(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, page)
}

// Create godoc
// @Summary Update an page
// @Description Update an new page item
// @Tags pages
// @Accept json
// @Produce json
// @Success 200 {object} Page
// @Failure 400 {object} support.HttpError
// @Router /api/pages [put]
func (controller pageController) Update(c echo.Context) error {

	page := &Page{}
	page.Name = c.FormValue("name")
	page.Title = c.FormValue("title")
	page.Meta = c.FormValue("meta")
	page.Content = c.FormValue("content")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, page)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "page created succesifully")
}

// @Summary Delte a page
// @Description Delte item
// @Tags pages
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Page
// @Failure 400 {object} support.HttpError
// @Router /api/pages [delete]
func (controller pageController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
