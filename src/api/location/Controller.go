package location

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"github.com/myrachanto/imagery"
)

// locationController ...
var (
	LocationController LocationControllerInterface = &locationController{}
)

type LocationControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	GetAll1(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type locationController struct {
	service LocationServiceInterface
}

func NewlocationController(ser LocationServiceInterface) LocationControllerInterface {
	return &locationController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a location
// @Description Create a new location item
// @Tags locations
// @Accept json
// @Produce json
// @Success 201 {object} location
// @Failure 400 {object} support.HttpError
// @Router /api/locations [post]

func (controller locationController) Create(c echo.Context) error {

	location := &Location{}
	// fmt.Println("--------------------create")
	location.Name = strings.TrimSpace(c.FormValue("name"))
	location.Title = c.FormValue("title")
	location.Description = c.FormValue("description")
	hasImage := c.FormValue("hasImage")
	location.Meta = c.FormValue("meta")
	location.Content = c.FormValue("content")
	location.PropertyType = c.FormValue("propertyType")

	url := strings.TrimSpace(c.FormValue("name"))
	location.Url = strings.Join(strings.Split(url, " "), "-")
	if hasImage == "yes" {
		pic, err2 := c.FormFile("picture")
		if err2 != nil {
			httperror := httperrors.NewBadRequestError("Invalid picture")
			return c.JSON(httperror.Code(), err2.Error())
		}
		filepath1, errs := controller.ProcessImage(pic)
		if errs != nil {
			return c.JSON(errs.Code(), errs.Message())
		}
		// fmt.Println("./////////////////step 5", filepath1)
		location.Picture = filepath1
	}
	_, err1 := controller.service.Create(location)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "location created succesifully")

}

// GetAll godoc
// @Summary GetAll a location
// @Description Getall locations
// @Tags locations
// @Accept json
// @Produce json
// @Success 201 {object} location
// @Failure 400 {object} support.HttpError
// @Router /api/locations [get]
func (controller locationController) GetAll(c echo.Context) error {
	// search := c.QueryParam("search")

	// fmt.Println("----------------------results")

	locations, err3 := controller.service.GetAll()

	// fmt.Println("----------------------results", locations)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, locations)
}

func (controller locationController) GetAll1(c echo.Context) error {

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
// @Summary Get a location
// @Description Get item
// @Tags locations
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} location
// @Failure 400 {object} support.HttpError
// @Router /api/locations [get]
func (controller locationController) GetOne(c echo.Context) error {
	code := c.Param("code")
	location, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, location)
}

// Create godoc
// @Summary Update an location
// @Description Update an new location item
// @Tags locations
// @Accept json
// @Produce json
// @Success 200 {object} location
// @Failure 400 {object} support.HttpError
// @Router /api/locations [put]
func (controller locationController) Update(c echo.Context) error {

	location := &Location{}
	location.Name = strings.TrimSpace(c.FormValue("name"))
	location.Title = c.FormValue("title")
	location.Description = c.FormValue("description")
	location.Meta = c.FormValue("meta")
	location.Content = c.FormValue("content")
	location.PropertyType = c.FormValue("propertyType")

	url := strings.TrimSpace(c.FormValue("name"))
	location.Url = strings.Join(strings.Split(url, " "), "-")
	hasImage := c.FormValue("hasImage")
	if hasImage == "yes" {
		pic, err2 := c.FormFile("picture")
		if err2 != nil {
			httperror := httperrors.NewBadRequestError("Invalid picture")
			return c.JSON(httperror.Code(), err2.Error())
		}
		filepath1, errs := controller.ProcessImage(pic)
		if errs != nil {
			return c.JSON(errs.Code(), errs.Message())
		}
		// fmt.Println("./////////////////step 5", filepath1)
		location.Picture = filepath1
	}
	code := c.Param("code")
	_, err1 := controller.service.Update(code, location)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "location created succesifully")
}

// Delete godoc
// @Summary Delte a location
// @Description Delte item
// @Tags locations
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} location
// @Failure 400 {object} support.HttpError
// @Router /api/locations [delete]
func (controller locationController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}

func (controller locationController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/locations/" + pic.Filename
	filePath1 := "/imgs/locations/" + pic.Filename
	// Destination
	dst, err4 := os.Create(filePath)
	if err4 != nil {
		return "", httperrors.NewBadRequestError("the Directory mess")
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		if err != nil {
			return "", httperrors.NewBadRequestError("error filling")
		}
	}
	imagery.Imageryrepository.Imagetype(filePath, filePath, 700, 1000)
	if _, err = io.Copy(dst, src); err != nil {
		if err != nil {
			return "", httperrors.NewBadRequestError("error saving the file")
		}
	}
	return filePath1, nil
}
