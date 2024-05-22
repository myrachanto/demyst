package category

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/estate/src/support"
)

// categoryController ...
var (
	CategoryController CategoryControllerInterface = &categoryController{}
)

type CategoryControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	GetAll1(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type categoryController struct {
	service CategoryServiceInterface
}

func NewcategoryController(ser CategoryServiceInterface) CategoryControllerInterface {
	return &categoryController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a category
// @Description Create a new category item
// @Tags categorys
// @Accept json
// @Produce json
// @Success 201 {object} Category
// @Failure 400 {object} support.HttpError
// @Router /api/categorys [post]

func (controller categoryController) Create(c echo.Context) error {

	category := &Category{}
	fmt.Println("--------------------create")
	category.Name = strings.TrimSpace(c.FormValue("name"))
	category.Title = c.FormValue("title")
	category.Description = c.FormValue("description")
	category.Meta = c.FormValue("meta")
	category.Content = c.FormValue("content")
	_, err1 := controller.service.Create(category)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "category created succesifully")

}

// GetAll godoc
// @Summary GetAll a category
// @Description Getall categorys
// @Tags categorys
// @Accept json
// @Produce json
// @Success 201 {object} Category
// @Failure 400 {object} support.HttpError
// @Router /api/categorys [get]
func (controller categoryController) GetAll(c echo.Context) error {
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
	categorys, err3 := controller.service.GetAll(searcher)

	// fmt.Println("----------------------results", categorys)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, categorys)
}
func (controller categoryController) GetAll1(c echo.Context) error {

	search := c.QueryParam("search")
	orderby := c.QueryParam("orderby")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh", ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	order, err := strconv.Atoi(c.QueryParam("order"))
	if err != nil {
		fmt.Println("Invalid pagesize")
		order = 1
	}
	searcher := support.Paginator{Page: page, Pagesize: pagesize, Search: search, Orderby: orderby, Order: order}
	// fmt.Println(">>>>>>>>>>>tag Bizname", Bizname)
	tags, err3 := controller.service.GetAll1(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, tags)
}

// GetOne godoc
// @Summary Get a category
// @Description Get item
// @Tags categorys
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Category
// @Failure 400 {object} support.HttpError
// @Router /api/categorys [get]
func (controller categoryController) GetOne(c echo.Context) error {
	code := c.Param("code")
	category, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, category)
}

// Create godoc
// @Summary Update an category
// @Description Update an new category item
// @Tags categorys
// @Accept json
// @Produce json
// @Success 200 {object} Category
// @Failure 400 {object} support.HttpError
// @Router /api/categorys [put]
func (controller categoryController) Update(c echo.Context) error {

	category := &Category{}
	category.Name = strings.TrimSpace(c.FormValue("name"))
	category.Title = c.FormValue("title")
	category.Description = c.FormValue("description")
	category.Meta = c.FormValue("meta")
	category.Content = c.FormValue("content")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, category)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "category created succesifully")
}

// Delete godoc
// @Summary Delte a category
// @Description Delte item
// @Tags categorys
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Category
// @Failure 400 {object} support.HttpError
// @Router /api/categorys [delete]
func (controller categoryController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
