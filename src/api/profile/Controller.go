package profile

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
	profilesections "github.com/myrachanto/estate/src/api/profileSections"
	"github.com/myrachanto/estate/src/support"
	"github.com/myrachanto/imagery"
)

// ProfileController ...
var (
	ProfileController ProfileControllerInterface = &profileController{}
)

type ProfileControllerInterface interface {
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

type profileController struct {
	service ProfileServiceInterface
}

func NewProfileController(ser ProfileServiceInterface) ProfileControllerInterface {
	return &profileController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a Profile
// @Description Create a new Profile item
// @Tags Profile
// @Accept json
// @Produce json
// @Success 201 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [post]
func (controller profileController) Create(c echo.Context) error {
	Profile := &Profile{}
	// fmt.Println("------------------step1")
	Profile.Name = c.FormValue("title")
	Profile.Title = c.FormValue("title")
	Profile.Url = c.FormValue("url")
	Profile.Meta = c.FormValue("meta")
	Profile.Content = c.FormValue("content")
	Profile.Caption = c.FormValue("caption")
	Profile.Sport = c.FormValue("category")
	Profile.Author = c.FormValue("author")
	Profile.Credit = c.FormValue("credit")
	Profile.PhotoCredit = c.FormValue("photocredit")
	featured, e := strconv.ParseBool(c.FormValue("featured"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, "could not parse featured")
	}
	Profile.Featured = featured
	exclusive, e := strconv.ParseBool(c.FormValue("exclusive"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, "could not parse exclusive")
	}
	Profile.Exclusive = exclusive
	trending, e := strconv.ParseBool(c.FormValue("trending"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, "could not parse featured")
	}
	Profile.Trending = trending

	fmt.Println("------------------step2")
	secs := []*profilesections.ProfileSection{}
	sections := c.FormValue("sections")
	if len(sections) > 0 {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(sections), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling products")
			return c.JSON(httperror.Code(), err4.Error())
		}
		// fmt.Println("./////////////////step4")

		for _, v := range producti {
			var sec profilesections.ProfileSection
			sec.Name = fmt.Sprintf("%s", v["name"])
			sec.Content = fmt.Sprintf("%s", v["content"])
			sec.Facebook = fmt.Sprintf("%s", v["facebook"])
			sec.Instagram = fmt.Sprintf("%s", v["instagram"])
			sec.Youtube = fmt.Sprintf("%s", v["youtube"])
			sec.Twitter = fmt.Sprintf("%s", v["twitter"])
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
	Profile.Picture = filepath1
	form, err := c.MultipartForm()

	files := form.File["pictures"]
	if err != nil {
		return err
	}
	if len(files) > 0 {
		for _, file := range files {
			filepath1, errs := controller.ProcessImage(file)
			if errs != nil {
				return c.JSON(errs.Code(), errs.Message())
			}
			sectnames := strings.Split(file.Filename, "_")
			sectname := sectnames[1]
			for _, g := range secs {
				if g.Name == sectname {
					g.Image = filepath1
				}
			}
		}
	}
	for _, g := range secs {
		Profile.Sections = append(Profile.Sections, g)
	}
	_, err1 := controller.service.Create(Profile)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "Profile created succesifully")

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
func (controller profileController) CreateComment(c echo.Context) error {

	com := &Comment{}
	fmt.Println("--------------------step 1")
	// com.Username = c.FormValue("name")
	// com.Email = c.FormValue("email")
	// com.Image = c.FormValue("image")
	// com.Code = c.FormValue("code")
	// com.Message = c.FormValue("message")
	err := c.Bind(&com)
	if err != nil {

		fmt.Println("--------------------step 2", err)
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println("--------------------step 3", com)
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
func (controller profileController) DeleteComment(c echo.Context) error {
	code := c.QueryParam("code")
	Profilecode := c.QueryParam("Profilecode")
	problem := controller.service.DeleteComment(code, Profilecode)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "comment successifuly deleted!")
}

// GetAll godoc
// @Summary GetAll a Profile
// @Description Getall Profile
// @Tags Profile
// @Accept json
// @Produce json
// @Success 201 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [get]
func (controller profileController) GetAll(c echo.Context) error {
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
	Profile, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, Profile)
}

// @Summary Get a Profile
// @Description Get item
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [get]
func (controller profileController) GetOne(c echo.Context) error {
	code := c.Param("code")
	Profile, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Profile)
}

// @Summary Get Profile by URL
// @Description Get Profile by URL
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /Profile [get]
func (controller profileController) GetOneByUrl(c echo.Context) error {
	code := c.Param("url")
	Profile, problem := controller.service.GetOneByUrl(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Profile)
}

// @Summary  Get  BY Category a Profile
// @Description  Get  BY Category item
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /Profile [get]
func (controller profileController) GetByCategory(c echo.Context) error {
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
	Profile, problem := controller.service.GetByCategory(code, searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, Profile)
}

// Create godoc
// @Summary Update an Profile
// @Description Update an new Profile item
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [put]
func (controller profileController) Update(c echo.Context) error {
	code := c.Param("code")
	Profile := &Profile{}
	// fmt.Println("------------------step1")
	Profile.Name = c.FormValue("title")
	Profile.Title = c.FormValue("title")
	Profile.Url = c.FormValue("url")
	Profile.Meta = c.FormValue("meta")
	Profile.Content = c.FormValue("content")
	Profile.Caption = c.FormValue("caption")
	Profile.Sport = c.FormValue("category")
	Profile.Author = c.FormValue("author")
	Profile.Credit = c.FormValue("credit")
	Profile.PhotoCredit = c.FormValue("photocredit")
	secs := []*profilesections.ProfileSection{}
	sections := c.FormValue("sections")

	fmt.Println("================>>>>>>++++++++++++++++ step 1", Profile.Credit)
	if string(sections) != "" {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(sections), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling products")
			return c.JSON(httperror.Code(), err4.Error())
		}
		for _, v := range producti {
			var sec profilesections.ProfileSection
			sec.Name = fmt.Sprintf("%s", v["name"])
			sec.Content = fmt.Sprintf("%s", v["content"])
			sec.Code = fmt.Sprintf("%s", v["code"])
			sec.Facebook = fmt.Sprintf("%s", v["facebook"])
			sec.Instagram = fmt.Sprintf("%s", v["instagram"])
			sec.Youtube = fmt.Sprintf("%s", v["youtube"])
			sec.Twitter = fmt.Sprintf("%s", v["twitter"])
			sec.Image = ""
			secs = append(secs, &sec)
		}
	}
	// fmt.Println("================>>>>>> step 2", sections)
	pic, err2 := c.FormFile("picture")
	if err2 == nil {
		filepath1, errs := controller.ProcessImage(pic)
		if errs != nil {
			return c.JSON(errs.Code(), errs.Message())
		}
		Profile.Picture = filepath1
	}
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["pictures"]

	// fmt.Println("================>>>>>> step 4", len(files))
	if len(files) > 0 {
		for _, file := range files {
			filepath1, errs := controller.ProcessImage(file)
			if errs != nil {
				return c.JSON(errs.Code(), errs.Message())
			}
			sectnames := strings.Split(file.Filename, "_")
			sectname := sectnames[1]
			for _, g := range secs {
				if g.Name == sectname {
					g.Image = filepath1
				}
			}
		}
	}
	for _, g := range secs {
		Profile.Sections = append(Profile.Sections, g)
	}
	fmt.Println("================>>>>>> step 5", Profile, code)
	_, err1 := controller.service.Update(code, Profile)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "Profile Updated succesifully")
}

// @Summary Update featured a Profile
// @Description Update featured item
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [delete]
func (controller profileController) UpdateFeatured(c echo.Context) error {
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

// @Summary Update trending a Profile
// @Description Update trending item
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [delete]
func (controller profileController) UpdateTrending(c echo.Context) error {
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

// @Summary Update exclusive a Profile
// @Description Update exclusive item
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [delete]
func (controller profileController) UpdateExclusive(c echo.Context) error {
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

// @Summary Delete a Profile
// @Description Delte item
// @Tags Profile
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Profile
// @Failure 400 {object} support.HttpError
// @Router /api/Profile [delete]
func (controller profileController) Delete(c echo.Context) error {
	id := string(c.Param("code"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller profileController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/profiles/" + pic.Filename
	filePath1 := "/imgs/profiles/" + pic.Filename
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
// 		filePath := "./src/public/imgs/Profile/" + pic.Filename
// 		filePath1 := "/imgs/Profile/" + pic.Filename
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

// 		Profile.Picture = filePath1
// 		_, err1 := controller.service.Create(Profile)
// 		if err1 != nil {
// 			return c.JSON(err1.Code(), err1)
// 		}
// 		if _, err = io.Copy(dst, src); err != nil {
// 			if err2 != nil {
// 				httperror := httperrors.NewBadRequestError("error filling")
// 				return c.JSON(httperror.Code(), httperror.Message())
// 			}
// 		}
// 		return c.JSON(http.StatusCreated, "Profile Updated succesifully")
