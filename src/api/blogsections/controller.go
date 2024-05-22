package BlogSections

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/imagery"
)

// blogsectionssectionsController ...
var (
	BlogSectionController BlogSectionControllerInterface = &blogSectionsController{}
)

type BlogSectionControllerInterface interface {
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type blogSectionsController struct {
	service BlogSectionServiceInterface
}

func NewblogSectionsController(ser BlogSectionServiceInterface) BlogSectionControllerInterface {
	return &blogSectionsController{
		ser,
	}
}

// Create godoc
// @Summary Update an blogsections
// @Description Update an new blogsections item
// @Tags blogsections
// @Accept json
// @Produce json
// @Success 200 {object} blog
// @Failure 400 {object} support.HttpError
// @Router /api/blogsections [put]
func (controller blogSectionsController) Update(c echo.Context) error {
	code := c.Param("code")
	blogsections := &BlogSection{}
	blogsections.Name = c.FormValue("name")
	blogsections.Content = c.FormValue("content")
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
	blogsections.Image = filepath1
	_, err1 := controller.service.Update(code, blogsections)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "blogsections Updated succesifully")
}

// @Summary Delete a blogsections
// @Description Delte item
// @Tags blogsections
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} blog
// @Failure 400 {object} support.HttpError
// @Router /api/blogsections [delete]
func (controller blogSectionsController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller blogSectionsController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/blog/" + pic.Filename
	filePath1 := "/imgs/blog/" + pic.Filename
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
