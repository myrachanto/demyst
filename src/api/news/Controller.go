package news

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/imagery"
	newssections "github.com/myrachanto/sports/src/api/newsSections"
	"github.com/myrachanto/sports/src/support"
)

// newsController ...
var (
	NewsController NewsControllerInterface = &newsController{}
)

type NewsControllerInterface interface {
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
}

type newsController struct {
	service NewsServiceInterface
}

func NewnewsController(ser NewsServiceInterface) NewsControllerInterface {
	return &newsController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a news
// @Description Create a new news item
// @Tags news
// @Accept json
// @Produce json
// @Success 201 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [post]
func (controller newsController) Create(c echo.Context) error {
	news := &News{}
	// fmt.Println("------------------step1")
	news.Name = c.FormValue("title")
	news.Title = c.FormValue("title")
	news.Url = c.FormValue("url")
	news.Meta = c.FormValue("meta")
	news.Content = c.FormValue("content")
	news.Caption = c.FormValue("caption")
	news.Sport = c.FormValue("category")
	news.Author = c.FormValue("author")
	news.Credit = c.FormValue("credit")
	featured, e := strconv.ParseBool(c.FormValue("featured"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, "could not parse featured")
	}
	news.Featured = featured
	exclusive, e := strconv.ParseBool(c.FormValue("exclusive"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, "could not parse exclusive")
	}
	news.Exclusive = exclusive
	trending, e := strconv.ParseBool(c.FormValue("trending"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, "could not parse featured")
	}
	news.Trending = trending

	// fmt.Println("------------------step2")
	secs := []*newssections.NewsSection{}
	sections := c.FormValue("sections")
	if string(sections) != "0" {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(sections), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling products")
			return c.JSON(httperror.Code(), err4.Error())
		}
		// fmt.Println("./////////////////step4")

		for _, v := range producti {
			var sec newssections.NewsSection
			sec.Name = fmt.Sprintf("%s", v["name"])
			sec.Content = fmt.Sprintf("%s", v["content"])
			sec.Image = ""
			secs = append(secs, &sec)
		}
		// fmt.Println("./////////////////", ts)
	}
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
	news.Picture = filepath1
	form, err := c.MultipartForm()

	files := form.File["pictures"]
	if err != nil {
		return err
	}
	if len(files) > 0 {
		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				httperror := httperrors.NewBadRequestError("Invalid picture")
				return c.JSON(httperror.Code(), err.Error())
			}
			defer src.Close()
			filePath := "./src/public/imgs/news/" + file.Filename
			dst, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer dst.Close()

			// fmt.Println("------------------step4", files)
			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			sectname := strings.Split(file.Filename, "_")[1]
			// Destination
			Original_Path := strings.Split(file.Filename, ".")
			name1 := Original_Path[len(Original_Path)-1]
			nameSplit := strings.Join(strings.Split(news.Name, " "), "-")
			// imagename := nameSplit + "." + name1

			imagename := fmt.Sprintf("%s.%s", nameSplit, name1)
			imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
			filepath3 := "./src/public/imgs/news/" + imagename
			support.RenameImage(filePath, filepath3)
			filepath5 := "/imgs/news/" + imagename
			for _, g := range secs {
				// fmt.Println("------------------step8", g.Name, g.Content, sectname)
				if g.Name == sectname {
					g.Image = filepath5
					// fmt.Println("------------------step8", g.Image)
				}
			}
		}
	}
	// fmt.Println("------------------step11", secs)
	// fmt.Println("./////////////////step5")
	for _, g := range secs {
		news.Sections = append(news.Sections, *g)
	}
	// news.Sections = secs
	// fmt.Println("------------------step6")
	_, err1 := controller.service.Create(news)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "news created succesifully")

}

// GetAll godoc
// @Summary GetAll a news
// @Description Getall news
// @Tags news
// @Accept json
// @Produce json
// @Success 201 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [get]
func (controller newsController) GetAll(c echo.Context) error {
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
	news, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, news)
}

// @Summary Get a news
// @Description Get item
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [get]
func (controller newsController) GetOne(c echo.Context) error {
	code := c.Param("code")
	news, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, news)
}

// @Summary Get news by URL
// @Description Get news by URL
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} News
// @Failure 400 {object} support.HttpError
// @Router /news [get]
func (controller newsController) GetOneByUrl(c echo.Context) error {
	code := c.Param("url")
	news, problem := controller.service.GetOneByUrl(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, news)
}

// @Summary  Get  BY Category a news
// @Description  Get  BY Category item
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} News
// @Failure 400 {object} support.HttpError
// @Router /news [get]
func (controller newsController) GetByCategory(c echo.Context) error {
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
	news, problem := controller.service.GetByCategory(code, searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, news)
}

// Create godoc
// @Summary Update an news
// @Description Update an new news item
// @Tags news
// @Accept json
// @Produce json
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [put]
func (controller newsController) Update(c echo.Context) error {
	code := c.Param("code")
	news := &News{}
	// fmt.Println("------------------step1")
	news.Name = c.FormValue("title")
	news.Title = c.FormValue("title")
	news.Url = c.FormValue("url")
	news.Meta = c.FormValue("meta")
	news.Content = c.FormValue("content")
	news.Caption = c.FormValue("caption")
	news.Sport = c.FormValue("category")
	news.Author = c.FormValue("author")
	news.Credit = c.FormValue("credit")
	// featured, e := strconv.ParseBool(c.FormValue("featured"))
	// if e != nil {
	// 	return c.JSON(http.StatusBadRequest, "could not parse featured")
	// }
	// news.Featured = featured
	// exclusive, e := strconv.ParseBool(c.FormValue("exclusive"))
	// if e != nil {
	// 	return c.JSON(http.StatusBadRequest, "could not parse exclusive")
	// }
	// news.Exclusive = exclusive
	// trending, e := strconv.ParseBool(c.FormValue("trending"))
	// if e != nil {
	// 	return c.JSON(http.StatusBadRequest, "could not parse featured")
	// }
	// news.Trending = trending

	// fmt.Println("------------------step2")
	secs := []*newssections.NewsSection{}
	sections := c.FormValue("sections")

	// fmt.Println("------------------step2a", sections)
	if string(sections) != "" {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(sections), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling products")
			return c.JSON(httperror.Code(), err4.Error())
		}
		// fmt.Println("./////////////////step4")

		for _, v := range producti {
			var sec newssections.NewsSection
			sec.Name = fmt.Sprintf("%s", v["name"])
			sec.Content = fmt.Sprintf("%s", v["content"])
			sec.Image = ""
			secs = append(secs, &sec)
		}
		// fmt.Println("./////////////////", ts)
	}
	// fmt.Println("------------------step3")
	pic, err2 := c.FormFile("picture")
	//    fmt.Println(pic.Filename)
	// fmt.Println("------------------step3a", err2)
	if err2 == nil {
		filepath1, errs := controller.ProcessImage(pic)
		if errs != nil {
			return c.JSON(errs.Code(), errs.Message())
		}
		news.Picture = filepath1
	}

	// fmt.Println("------------------step4")
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["pictures"]

	// fmt.Println("------------------step5")
	if len(files) > 0 {
		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				httperror := httperrors.NewBadRequestError("Invalid picture")
				return c.JSON(httperror.Code(), err.Error())
			}
			defer src.Close()
			filePath := "./src/public/imgs/news/" + file.Filename
			dst, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer dst.Close()

			// fmt.Println("------------------step4", files)
			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			sectname := strings.Split(file.Filename, "_")[1]
			// Destination
			Original_Path := strings.Split(file.Filename, ".")
			name1 := Original_Path[len(Original_Path)-1]
			nameSplit := strings.Join(strings.Split(news.Name, " "), "-")
			// imagename := nameSplit + "." + name1

			imagename := fmt.Sprintf("%s.%s", nameSplit, name1)
			imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
			filepath3 := "./src/public/imgs/news/" + imagename
			support.RenameImage(filePath, filepath3)
			filepath5 := "/imgs/news/" + imagename
			for _, g := range secs {
				// fmt.Println("------------------step8", g.Name, g.Content, sectname)
				if g.Name == sectname {
					g.Image = filepath5
					// fmt.Println("------------------step8", g.Image)
				}
			}
		}
	}
	// fmt.Println("------------------step11", secs)
	// fmt.Println("./////////////////step5")
	for _, g := range secs {
		news.Sections = append(news.Sections, *g)
	}
	_, err1 := controller.service.Update(code, news)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "news Updated succesifully")
}

// @Summary Update featured a news
// @Description Update featured item
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [delete]
func (controller newsController) UpdateFeatured(c echo.Context) error {
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

// @Summary Update trending a news
// @Description Update trending item
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [delete]
func (controller newsController) UpdateTrending(c echo.Context) error {
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

// @Summary Update exclusive a news
// @Description Update exclusive item
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [delete]
func (controller newsController) UpdateExclusive(c echo.Context) error {
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

// @Summary Delete a news
// @Description Delte item
// @Tags news
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} News
// @Failure 400 {object} support.HttpError
// @Router /api/news [delete]
func (controller newsController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller newsController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
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
// 		filePath := "./src/public/imgs/news/" + pic.Filename
// 		filePath1 := "/imgs/news/" + pic.Filename
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

// 		news.Picture = filePath1
// 		_, err1 := controller.service.Create(news)
// 		if err1 != nil {
// 			return c.JSON(err1.Code(), err1)
// 		}
// 		if _, err = io.Copy(dst, src); err != nil {
// 			if err2 != nil {
// 				httperror := httperrors.NewBadRequestError("error filling")
// 				return c.JSON(httperror.Code(), httperror.Message())
// 			}
// 		}
// 		return c.JSON(http.StatusCreated, "news Updated succesifully")
