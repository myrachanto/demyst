package newssections

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/imagery"
)

// newssectionssectionsController ...
var (
	NewsSectionController NewsSectionControllerInterface = &newsSectionsController{}
)

type NewsSectionControllerInterface interface {
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type newsSectionsController struct {
	service NewsSectionServiceInterface
}

func NewnewsSectionsController(ser NewsSectionServiceInterface) NewsSectionControllerInterface {
	return &newsSectionsController{
		ser,
	}
}

// Create godoc
// @Summary Update an newssections
// @Description Update an new newssections item
// @Tags newssections
// @Accept json
// @Produce json
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/newssections [put]
func (controller newsSectionsController) Update(c echo.Context) error {
	code := c.Param("code")
	newssections := &NewsSection{}
	newssections.Name = c.FormValue("name")
	newssections.Content = c.FormValue("content")
	pic, err2 := c.FormFile("picture")
	//    fmt.Println(pic.Filename)
	if err2 != nil {
		httperror := httperrors.NewBadRequestError("Invalid picture")
		return c.JSON(httperror.Code(), err2.Error())
	}
	filepath1, errs := controller.ProcessImage(pic)
	if errs != nil {
		return c.JSON(errs.Code(), errs.Message())
	}
	newssections.Image = filepath1
	_, err1 := controller.service.Update(code, newssections)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "newssections Updated succesifully")
}

// @Summary Delete a newssections
// @Description Delte item
// @Tags newssections
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/newssections [delete]
func (controller newsSectionsController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller newsSectionsController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/news/" + pic.Filename
	filePath1 := "/imgs/news/" + pic.Filename
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
	imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
	if _, err = io.Copy(dst, src); err != nil {
		if err != nil {
			return "", httperrors.NewBadRequestError("error saving the file")
		}
	}
	return filePath1, nil
}
