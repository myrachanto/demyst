package blog

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"

	// blogsections "github.com/myrachanto/estate/src/api/BlogSections"
	"github.com/myrachanto/estate/src/support"
	"github.com/myrachanto/imagery"
)

// BlogController ...
var (
	BlogController BlogControllerInterface = &blogController{}
)

type BlogControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	UpdateFeatured(c echo.Context) error
	UpdateTrending(c echo.Context) error
	UpdateExclusive(c echo.Context) error
	Delete(c echo.Context) error
	GetOneByUrl(c echo.Context) error
	GetByCategory(c echo.Context) error
	CreateComment(c echo.Context) error
	DeleteComment(c echo.Context) error
}

type blogController struct {
	service BlogServiceInterface
}

func NewBlogController(ser BlogServiceInterface) BlogControllerInterface {
	return &blogController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a Blog
// @Description Create a new Blog item
// @Tags Blog
// @Accept json
// @Produce json
// @Success 201 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [post]
func (controller blogController) Create(c echo.Context) error {
	Blog := &Blog{}
	// fmt.Println("------------------step1")
	Blog.Name = c.FormValue("title")
	Blog.Title = c.FormValue("title")
	Blog.Url = c.FormValue("url")
	Blog.Meta = c.FormValue("meta")
	Blog.Content = c.FormValue("content")
	Blog.Caption = c.FormValue("caption")
	hasImage := c.FormValue("hasImage")
	if hasImage == "yes" {
		pic, err2 := c.FormFile("picture")
		fmt.Println(pic.Filename)
		if err2 != nil {
			httperror := httperrors.NewBadRequestError("Invalid picture")
			return c.JSON(httperror.Code(), err2.Error())
		}
		filepath1, errs := controller.ProcessImage(pic)
		if errs != nil {
			return c.JSON(errs.Code(), errs.Message())
		}
		// fmt.Println("./////////////////step 5", filepath1)
		Blog.Picture = filepath1
	}
	_, err1 := controller.service.Create(Blog)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "Blog created succesifully")

}

// Create Comment godoc
// @Summary Create a Comment
// @Description Create a new Comment item
// @Tags Comments
// @Accept json
// @Produce json
// @Success 201 {object} Comment
// @Failure 400 {object} support.HttpError
// @Router /api/comment/create [post]
func (controller blogController) CreateComment(c echo.Context) error {

	com := &Comment{}
	fmt.Println("--------------------step 1")
	// com.Username = c.FormValue("name")
	// com.Email = c.FormValue("email")
	// com.Image = c.FormValue("image")
	// com.Code = c.FormValue("code")
	// com.Message = c.FormValue("message")
	err := c.Bind(&com)
	if err != nil {

		return c.String(http.StatusBadRequest, "bad request")
	}
	err1 := controller.service.CreateComment(com)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "page created succesifully")

}

// Delete Comment godoc
// @Summary Delete a Comment
// @Description Delete a new Comment item
// @Tags Comments
// @Accept json
// @Produce json
// @Success 201 {object} Comment
// @Failure 400 {object} support.HttpError
// @Router /api/comment/delete [post]
func (controller blogController) DeleteComment(c echo.Context) error {
	code := c.QueryParam("code")
	Blogcode := c.QueryParam("Blogcode")
	problem := controller.service.DeleteComment(code, Blogcode)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "comment successifuly deleted!")
}

// GetAll godoc
// @Summary GetAll a Blog
// @Description Getall Blog
// @Tags Blog
// @Accept json
// @Produce json
// @Success 201 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [get]
func (controller blogController) GetAll(c echo.Context) error {
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
	Blog, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, Blog)
}

// @Summary Get a Blog
// @Description Get item
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [get]
func (controller blogController) GetOne(c echo.Context) error {
	code := c.Param("code")
	Blog, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Blog)
}

// @Summary Get Blog by URL
// @Description Get Blog by URL
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /Blog [get]
func (controller blogController) GetOneByUrl(c echo.Context) error {
	code := c.Param("url")
	Blog, problem := controller.service.GetOneByUrl(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Blog)
} 

// @Summary  Get  BY Category a Blog
// @Description  Get  BY Category item
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /Blog [get]
func (controller blogController) GetByCategory(c echo.Context) error {
	code := c.Param("category")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	fmt.Println("----------------------sdfgghh", code)
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator{Page: page, Pagesize: pagesize, Search: code}
	Blog, problem := controller.service.GetByCategory(code, searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Blog)
}

// Create godoc
// @Summary Update an Blog
// @Description Update an new Blog item
// @Tags Blog
// @Accept json
// @Produce json
// @Success 200 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [put]
func (controller blogController) Update(c echo.Context) error {
	code := c.Param("code")
	Blog := &Blog{}
	// fmt.Println("------------------step1")
	Blog.Name = c.FormValue("title")
	Blog.Title = c.FormValue("title")
	Blog.Url = c.FormValue("url")
	Blog.Meta = c.FormValue("meta")
	Blog.Content = c.FormValue("content")
	hasImage := c.FormValue("hasImage")
	if hasImage == "yes" {
		pic, err2 := c.FormFile("picture")
		if err2 == nil {
			filepath1, errs := controller.ProcessImage(pic)
			if errs != nil {
				return c.JSON(errs.Code(), errs.Message())
			}
			Blog.Picture = filepath1
		}
	}
	_, err1 := controller.service.Update(code, Blog)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "Blog Updated succesifully")
}

// @Summary Update featured a Blog
// @Description Update featured item
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [delete]
func (controller blogController) UpdateFeatured(c echo.Context) error {
	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>", code, feat)
	problem := controller.service.UpdateFeatured(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// @Summary Update trending a Blog
// @Description Update trending item
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [delete]
func (controller blogController) UpdateTrending(c echo.Context) error {
	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	problem := controller.service.UpdateTrending(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// @Summary Update exclusive a Blog
// @Description Update exclusive item
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [delete]
func (controller blogController) UpdateExclusive(c echo.Context) error {
	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	problem := controller.service.UpdateExclusive(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// @Summary Delete a Blog
// @Description Delte item
// @Tags Blog
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Blog
// @Failure 400 {object} support.HttpError
// @Router /api/Blog [delete]
func (controller blogController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller blogController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/blogs/" + pic.Filename
	filePath1 := "/imgs/blogs/" + pic.Filename
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

// pic, err2 := c.FormFile("picture")
// 	if pic != nil {
// 		//    fmt.Println(pic.Filename)
// 		if err2 != nil {
// 			httperror := httperrors.NewBadRequestError("Invalid picture")
// 			return c.JSON(httperror.Code(), err2.Error())
// 		}
// 		src, err := pic.Open()
// 		if err != nil {
// 			httperror := httperrors.NewBadRequestError("the picture is corrupted")
// 			return c.JSON(httperror.Code(), err.Error())
// 		}
// 		defer src.Close()
// 		// filePath := "./public/imgs/blogs/"
// 		filePath := "./src/public/imgs/Blog/" + pic.Filename
// 		filePath1 := "/imgs/Blog/" + pic.Filename
// 		// Destination
// 		dst, err4 := os.Create(filePath)
// 		if err4 != nil {
// 			httperror := httperrors.NewBadRequestError("the Directory mess")
// 			return c.JSON(httperror.Code(), err4.Error())
// 		}
// 		defer dst.Close()
// 		// Copy
// 		if _, err = io.Copy(dst, src); err != nil {
// 			if err2 != nil {
// 				httperror := httperrors.NewBadRequestError("error filling")
// 				return c.JSON(httperror.Code(), httperror.Message())
// 			}
// 		}
// 		imagery.Imageryrepository.Imagetype(filePath, filePath, 300, 500)

// 		Blog.Picture = filePath1
// 		_, err1 := controller.service.Create(Blog)
// 		if err1 != nil {
// 			return c.JSON(err1.Code(), err1)
// 		}
// 		if _, err = io.Copy(dst, src); err != nil {
// 			if err2 != nil {
// 				httperror := httperrors.NewBadRequestError("error filling")
// 				return c.JSON(httperror.Code(), httperror.Message())
// 			}
// 		}
// 		return c.JSON(http.StatusCreated, "Blog Updated succesifully")
