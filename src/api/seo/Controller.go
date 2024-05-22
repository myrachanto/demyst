package seo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/estate/src/support"
)

// seoController ...
var (
	SeoController SeoControllerInterface = &seoController{}
)

type SeoControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	GetAll1(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type seoController struct {
	service SeoServiceInterface
}

func NewseoController(ser SeoServiceInterface) SeoControllerInterface {
	return &seoController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a seo
// @Description Create a new seo item
// @Tags seos
// @Accept json
// @Produce json
// @Success 201 {object} seo
// @Failure 400 {object} support.HttpError
// @Router /api/seos [post]

func (controller seoController) Create(c echo.Context) error {

	seo := &Seo{}
	// fmt.Println("--------------------create")
	seo.Page = c.FormValue("page")
	seo.Name = c.FormValue("name")
	seo.Title = c.FormValue("title")
	seo.Description = c.FormValue("description")
	seo.Meta = c.FormValue("meta")
	seo.Kind = c.FormValue("kind")
	seo.Location = c.FormValue("location")
	seo.Sublocation = c.FormValue("sublocation")
	_, err1 := controller.service.Create(seo)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "seo created succesifully")

}

// GetAll godoc
// @Summary GetAll a seo
// @Description Getall seos
// @Tags seos
// @Accept json
// @Produce json
// @Success 201 {object} seo
// @Failure 400 {object} support.HttpError
// @Router /api/seos [get]
func (controller seoController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")

	fmt.Println("----------------------results")
	seos, err3 := controller.service.GetAll(search)

	// fmt.Println("----------------------results", seos)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, seos)
}
func (controller seoController) GetAll1(c echo.Context) error {

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
// @Summary Get a seo
// @Description Get item
// @Tags seos
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} seo
// @Failure 400 {object} support.HttpError
// @Router /api/seos [get]
func (controller seoController) GetOne(c echo.Context) error {
	code := c.Param("code")
	seo, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, seo)
}

// Create godoc
// @Summary Update an seo
// @Description Update an new seo item
// @Tags seos
// @Accept json
// @Produce json
// @Success 200 {object} seo
// @Failure 400 {object} support.HttpError
// @Router /api/seos [put]
func (controller seoController) Update(c echo.Context) error {

	seo := &Seo{}
	seo.Name = c.FormValue("name")
	seo.Title = c.FormValue("title")
	seo.Page = c.FormValue("page")
	seo.Description = c.FormValue("description")
	seo.Meta = c.FormValue("meta")
	seo.Kind = c.FormValue("kind")
	seo.Location = c.FormValue("location")
	seo.Sublocation = c.FormValue("sublocation")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, seo)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "seo created succesifully")
}

// Delete godoc
// @Summary Delte a seo
// @Description Delte item
// @Tags seos
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} seo
// @Failure 400 {object} support.HttpError
// @Router /api/seos [delete]
func (controller seoController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
