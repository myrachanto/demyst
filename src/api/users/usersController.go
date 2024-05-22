package users

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"github.com/myrachanto/imagery"
)

// UserController ...
var (
	UserController UserControllerInterface = &userController{}
)

type UserControllerInterface interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	// RenewAccessToken(c echo.Context) error
	Logout(c echo.Context) error
	GetOne(c echo.Context) error
	Forgot(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	PasswordUpdate(c echo.Context) error
	PasswordReset(c echo.Context) error
	Delete(c echo.Context) error
	UpdateAdmin(c echo.Context) error
	UpdateAuditor(c echo.Context) error
}

type userController struct {
	service UserServiceInterface
}

func NewUserController(ser UserServiceInterface) UserControllerInterface {
	return &userController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a user
// @Description Create a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /users [post]
func (controller userController) Create(c echo.Context) error {

	user := &User{}
	user.Fullname = c.FormValue("fullname")
	user.Username = c.FormValue("username")
	user.Phone = c.FormValue("phone")
	// user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	auth, err1 := controller.service.Create(user)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, auth)
}

// Login godoc
// @Summary Login a user
// @Description Login user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /login [post]
func (controller userController) Login(c echo.Context) error {
	user := &LoginUser{}
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	// user.UserAgent = c.Request().UserAgent()
	auth, problem := controller.service.Login(user)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, auth)
}

// // renewAccesstoken godoc
// // @Summary renewAccesstoken a user
// // @Description renewAccesstoken user
// // @Tags users
// // @Accept json
// // @Produce json
// // @Success 200 {object} User
// // @Failure 400 {object} support.HttpError
// // @Router /front/renewAccesstoken [post]
// func (controller userController) RenewAccessToken(c echo.Context) error {
// 	renewaccessToken := c.FormValue("renewaccessToken")
// 	// fmt.Println(">>>>>>>>>>>>>>>>>>>", renewaccessToken)
// 	auth, problem := controller.service.RenewAccessToken(renewaccessToken)
// 	if problem != nil {
// 		fmt.Println(problem)
// 		return c.JSON(problem.Code(), problem.Message())
// 	}
// 	return c.JSON(http.StatusOK, auth)
// }

// logout godoc
// @Summary logout a user
// @Description logout user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/logout [post]
func (controller userController) Logout(c echo.Context) error {
	token := string(c.Param("token"))
	_, problem := controller.service.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")
}

// GetAll godoc
// @Summary GetAll a user
// @Description Getall users
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [get]
func (controller userController) GetAll(c echo.Context) error {
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

	users, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, users)
}

// @Summary Get a user
// @Description Get item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [get]
func (controller userController) GetOne(c echo.Context) error {
	code := c.Param("code")
	user, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, user)
}

// Forgot godoc
// @Summary Forgot a user
// @Description Forgot user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /forgot [post]
func (controller userController) Forgot(c echo.Context) error {
	email := c.FormValue("email")
	_, problem := controller.service.Forgot(email)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// Forgot godoc
// @Summary Forgot a user
// @Description Forgot user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /front/forgot [post]
func (controller userController) PasswordUpdate(c echo.Context) error {
	// fmt.Println("-----------------0")
	oldpassword := c.FormValue("oldpassword")
	email := c.FormValue("email")
	newpassword := c.FormValue("newpassword")
	_, _, problem := controller.service.PasswordUpdate(oldpassword, email, newpassword)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// Reset godoc
// @Summary Reset a user
// @Description Reset user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users/reset [post]
func (controller userController) PasswordReset(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	_, problem := controller.service.PasswordReset(email, password)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update a user
// @Description Update a user item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [put]
func (controller userController) Update(c echo.Context) error {
	user := &User{}
	user.Fullname = c.FormValue("fullname")
	user.Username = c.FormValue("username")
	user.Phone = c.FormValue("phone")
	code := c.Param("code")
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
	user.Picture = filepath1
	problem := controller.service.Update(code, user)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusCreated, "Updated successifuly")
}

// Update admin godoc
// @Summary Update a user Admin
// @Description Update a user item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/admin/update/:code [put]
func (controller userController) UpdateAdmin(c echo.Context) error {

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>> step1")
	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>", code, feat)
	problem := controller.service.UpdateAdmin(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// Update Auditor godoc
// @Summary Update a user Auditor
// @Description Update a user Auditor item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/admin/update/:code [put]
func (controller userController) UpdateAuditor(c echo.Context) error {

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>> step1")
	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>", code, feat)
	problem := controller.service.UpdateAuditor(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// Delete godoc
// @Summary Delete a user
// @Description Create a new user item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} string
// @Failure 400 {object} support.HttpError
// @Router /api/users [delete]
func (controller userController) Delete(c echo.Context) error {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
func (controller userController) ProcessImage(pic *multipart.FileHeader) (string, httperrors.HttpErr) {
	src, err := pic.Open()
	if err != nil {
		return "", httperrors.NewBadRequestError("the picture is corrupted")
	}
	defer src.Close()
	// filePath := "./public/imgs/blogs/"
	filePath := "./src/public/imgs/users/" + pic.Filename
	filePath1 := "/imgs/users/" + pic.Filename
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
	imagery.Imageryrepository.Imagetype(filePath, filePath, 300, 500)
	if _, err = io.Copy(dst, src); err != nil {
		if err != nil {
			return "", httperrors.NewBadRequestError("error saving the file")
		}
	}
	return filePath1, nil
}
