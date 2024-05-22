package feature

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/estate/src/support"
)

// featureController ...
var (
	FeatureController FeatureControllerInterface = &featureController{}
)

type FeatureControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	GetAll1(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type featureController struct {
	service FeatureServiceInterface
}

func NewfeatureController(ser FeatureServiceInterface) FeatureControllerInterface {
	return &featureController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a feature
// @Description Create a new feature item
// @Tags features
// @Accept json
// @Produce json
// @Success 201 {object} feature
// @Failure 400 {object} support.HttpError
// @Router /api/features [post]

func (controller featureController) Create(c echo.Context) error {

	feature := &Feature{}
	// fmt.Println("--------------------create")
	feature.Name = strings.TrimSpace(c.FormValue("name"))
	feature.Title = c.FormValue("title")
	feature.Description = c.FormValue("description")
	_, err1 := controller.service.Create(feature)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "feature created succesifully")

}

// GetAll godoc
// @Summary GetAll a feature
// @Description Getall features
// @Tags features
// @Accept json
// @Produce json
// @Success 201 {object} feature
// @Failure 400 {object} support.HttpError
// @Router /api/features [get]
func (controller featureController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")

	fmt.Println("----------------------results")
	features, err3 := controller.service.GetAll(search)

	// fmt.Println("----------------------results", features)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, features)
}
func (controller featureController) GetAll1(c echo.Context) error {

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
// @Summary Get a feature
// @Description Get item
// @Tags features
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} feature
// @Failure 400 {object} support.HttpError
// @Router /api/features [get]
func (controller featureController) GetOne(c echo.Context) error {
	code := c.Param("code")
	feature, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, feature)
}

// Create godoc
// @Summary Update an feature
// @Description Update an new feature item
// @Tags features
// @Accept json
// @Produce json
// @Success 200 {object} feature
// @Failure 400 {object} support.HttpError
// @Router /api/features [put]
func (controller featureController) Update(c echo.Context) error {

	feature := &Feature{}
	feature.Name = strings.TrimSpace(c.FormValue("name"))
	feature.Title = c.FormValue("title")
	feature.Description = c.FormValue("description")
	code := c.Param("code")
	_, err1 := controller.service.Update(code, feature)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "feature created succesifully")
}

// Delete godoc
// @Summary Delte a feature
// @Description Delte item
// @Tags features
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} feature
// @Failure 400 {object} support.HttpError
// @Router /api/features [delete]
func (controller featureController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
