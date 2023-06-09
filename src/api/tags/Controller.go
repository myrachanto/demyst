package tags

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/sports/src/support"
)

// tagController ...
var (
	TagController TagControllerInterface = &tagController{}
)

type TagControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type tagController struct {
	service TagServiceInterface
}

func NewtagController(ser TagServiceInterface) TagControllerInterface {
	return &tagController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a tag
// @Description Create a new tag item
// @Tags tags
// @Accept json
// @Produce json
// @Success 201 {object} Tag
// @Failure 400 {object} support.HttpError
// @Router /api/tags [post]

func (controller tagController) Create(c echo.Context) error {

	tag := &Tag{}
	// fmt.Println("--------------------create")
	tag.Name = c.FormValue("name")
	tag.Title = c.FormValue("title")
	tag.Description = c.FormValue("description")
	_, err1 := controller.service.Create(tag)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "tag created succesifully")

}

// GetAll godoc
// @Summary GetAll a tag
// @Description Getall tags
// @Tags tags
// @Accept json
// @Produce json
// @Success 201 {object} Tag
// @Failure 400 {object} support.HttpError
// @Router /api/tags [get]
func (controller tagController) GetAll(c echo.Context) error {
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
	tags, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, tags)
}

// @Summary Get a tag
// @Description Get item
// @Tags tags
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Tag
// @Failure 400 {object} support.HttpError
// @Router /api/tags [get]
func (controller tagController) GetOne(c echo.Context) error {
	code := c.Param("code")
	tag, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, tag)
}

// Create godoc
// @Summary Update an tag
// @Description Update an new tag item
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} Tag
// @Failure 400 {object} support.HttpError
// @Router /api/tags [put]
func (controller tagController) Update(c echo.Context) error {

	tag := &Tag{}
	tag.Name = c.FormValue("name")
	tag.Title = c.FormValue("title")
	tag.Description = c.FormValue("description")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, tag)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "tag created succesifully")
}

// @Summary Delte a tag
// @Description Delte item
// @Tags tags
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Tag
// @Failure 400 {object} support.HttpError
// @Router /api/tags [delete]
func (controller tagController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
