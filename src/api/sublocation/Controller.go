package subLocation

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/estate/src/support"
)

// subLocationController ...
var (
	SubLocationController SubLocationControllerInterface = &subLocationController{}
)

type SubLocationControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	GetAll1(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetAllByLocation(c echo.Context) error
}

type subLocationController struct {
	service SubLocationServiceInterface
}

func NewsubLocationController(ser SubLocationServiceInterface) SubLocationControllerInterface {
	return &subLocationController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a subLocation
// @Description Create a new subLocation item
// @Tags subLocations
// @Accept json
// @Produce json
// @Success 201 {object} subLocation
// @Failure 400 {object} support.HttpError
// @Router /api/subLocations [post]

func (controller subLocationController) Create(c echo.Context) error {

	subLocation := &SubLocation{}
	// fmt.Println("--------------------create")
	subLocation.Name = strings.TrimSpace(c.FormValue("name"))
	subLocation.Title = c.FormValue("title")
	subLocation.Description = c.FormValue("description")
	subLocation.Location = c.FormValue("location")
	subLocation.Meta = c.FormValue("meta")
	subLocation.Content = c.FormValue("content")
	subLocation.PropertyType = c.FormValue("propertyType")
	url := strings.TrimSpace(c.FormValue("name"))
	subLocation.Url = strings.Join(strings.Split(url, " "), "-")
	_, err1 := controller.service.Create(subLocation)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "subLocation created succesifully")

}

// GetAll godoc
// @Summary GetAll a subLocation
// @Description Getall subLocations
// @Tags subLocations
// @Accept json
// @Produce json
// @Success 201 {object} subLocation
// @Failure 400 {object} support.HttpError
// @Router /api/subLocations [get]
func (controller subLocationController) GetAll(c echo.Context) error {
	subLocations, err3 := controller.service.GetAll()

	// fmt.Println("----------------------results", subLocations)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, subLocations)
}
func (controller subLocationController) GetAll1(c echo.Context) error {

	search := c.QueryParam("search")
	orderby := c.QueryParam("orderby")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	fmt.Println("----------------------sdfgghh", ps)
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
	tags, err3 := controller.service.GetAll1(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	fmt.Println(">>>>>>>>>>>tag", tags)
	return c.JSON(http.StatusOK, tags)
}

// GetOne godoc
// @Summary Get byLocation
// @Description Get item
// @Tags subLocations
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} subLocation
// @Failure 400 {object} support.HttpError
// @Router /api/subLocations [get]
func (controller subLocationController) GetAllByLocation(c echo.Context) error {
	code := c.Param("code")
	// fmt.Println("++++++++++++++++++", code)
	subLocations, problem := controller.service.GetAllByLocation(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, subLocations)
}

// GetOne godoc
// @Summary Get a subLocation
// @Description Get item
// @Tags subLocations
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} subLocation
// @Failure 400 {object} support.HttpError
// @Router /api/subLocations [get]
func (controller subLocationController) GetOne(c echo.Context) error {
	code := c.Param("code")
	subLocation, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, subLocation)
}

// Create godoc
// @Summary Update an subLocation
// @Description Update an new subLocation item
// @Tags subLocations
// @Accept json
// @Produce json
// @Success 200 {object} subLocation
// @Failure 400 {object} support.HttpError
// @Router /api/subLocations [put]
func (controller subLocationController) Update(c echo.Context) error {

	subLocation := &SubLocation{}
	subLocation.Name = strings.TrimSpace(c.FormValue("name"))
	subLocation.Title = c.FormValue("title")
	subLocation.Description = c.FormValue("description")
	subLocation.Location = c.FormValue("location")
	subLocation.Meta = c.FormValue("meta")
	subLocation.PropertyType = c.FormValue("propertyType")
	subLocation.Content = c.FormValue("content")
	url := strings.TrimSpace(c.FormValue("name"))
	subLocation.Url = strings.Join(strings.Split(url, " "), "-")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, subLocation)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "subLocation created succesifully")
}

// Delete godoc
// @Summary Delte a subLocation
// @Description Delte item
// @Tags subLocations
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} subLocation
// @Failure 400 {object} support.HttpError
// @Router /api/subLocations [delete]
func (controller subLocationController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
