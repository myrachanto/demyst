package majorcategory

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

// majorcategoryController ...
var (
	MajorcategoryController MajorcategoryControllerInterface = majorcategoryController{}
	Bizname                 string
)

type MajorcategoryControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	GetAll1(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type majorcategoryController struct {
	service MajorcategoryServiceInterface
}

func NewmajorcategoryController(ser MajorcategoryServiceInterface) MajorcategoryControllerInterface {
	return &majorcategoryController{
		ser,
	}
}

// ///////controllers/////////////////
// Create godoc
// @Summary Create a majorcategory
// @Description Create a new majorcategory item
// @Tags majorcategorys
// @Accept json
// @Produce json
// @Success 201 {object} Majorcategory
// @Failure 400 {object} support.HttpError
// @Router /api/majorcategorys [post]
func (controller majorcategoryController) Create(c echo.Context) error {

	majorcategory := &Majorcategory{}
	// code := c.Param("majorcode")
	majorcategory.Name = strings.TrimSpace(c.FormValue("name"))
	majorcategory.Description = c.FormValue("description")
	majorcategory.Title = c.FormValue("title")
	majorcategory.Url = c.FormValue("url")
	majorcategory.Meta = c.FormValue("meta")
	majorcategory.Content = c.FormValue("content")
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
		majorcategory.Picture = filepath1
	}

	_, err1 := controller.service.Create(majorcategory)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "created successifuly")
}

// GetAll godoc
// @Summary GetAll a majorcategory
// @Description Getall majorcategorys
// @Tags majorcategorys
// @Accept json
// @Produce json
// @Success 201 {object} Majorcategory
// @Failure 400 {object} support.HttpError
// @Router /api/majorcategorys [get]
func (controller majorcategoryController) GetAll(c echo.Context) error {

	if Bizname == "" {
		Bizname = c.QueryParam("bizname")
	}
	code := c.Param("majorcode")
	fmt.Println(code)
	majorcategorys, err3 := controller.service.GetAll()
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, majorcategorys)
}
func (controller majorcategoryController) GetAll1(c echo.Context) error {

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

// @Summary Get a majorcategory
// @Description Get item
// @Tags majorcategorys
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Majorcategory
// @Failure 400 {object} support.HttpError
// @Router /api/majorcategorys [get]
func (controller majorcategoryController) GetOne(c echo.Context) error {

	id := c.Param("code")
	majorcategory, problem := controller.service.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, majorcategory)
}

// GetOne godoc
// @Summary Update a majorcategory
// @Description Update a majorcategory item
// @Tags majorcategorys
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Majorcategory
// @Failure 400 {object} support.HttpError
// @Router /api/majorcategorys [put]
func (controller majorcategoryController) Update(c echo.Context) error {

	majorcategory := &Majorcategory{}
	majorcategory.Name = strings.TrimSpace(c.FormValue("name"))
	majorcategory.Description = c.FormValue("description")
	majorcategory.Title = c.FormValue("title")
	majorcategory.Url = c.FormValue("url")
	majorcategory.Meta = c.FormValue("meta")
	majorcategory.Content = c.FormValue("content")
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
		fmt.Println("./////////////////step 5", filepath1)
		majorcategory.Picture = filepath1
	}
	id := c.Param("code")
	_, problem := controller.service.Update(id, majorcategory)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusCreated, "updated successifuly")
}

// Delete godoc
// @Summary Delete a majorcategory
// @Description Create a new majorcategory item
// @Tags majorcategorys
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} string
// @Failure 400 {object} support.HttpError
// @Router /api/majorcategorys [delete]
func (controller majorcategoryController) Delete(c echo.Context) error {

	id := c.Param("code")
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller majorcategoryController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/properties/" + pic.Filename
	filePath1 := "/imgs/properties/" + pic.Filename
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
	imagery.Imageryrepository.Imagetype(filePath, filePath, 800, 600)
	if _, err = io.Copy(dst, src); err != nil {
		if err != nil {
			return "", httperrors.NewBadRequestError("error saving the file")
		}
	}
	return filePath1, nil
}
