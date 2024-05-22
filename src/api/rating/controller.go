package rating

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ratingController ...
var (
	RatingController RatingControllerInterface = ratingController{}
	Bizname          string
)

type RatingControllerInterface interface {
	Create(c echo.Context) error
	Create1(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Featured(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type ratingController struct {
	service RatingServiceInterface
}

func NewratingController(ser RatingServiceInterface) RatingControllerInterface {
	return &ratingController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a rating
// @Description Create a new rating item
// @Tags ratings
// @Accept json
// @Produce json
// @Success 201 {object} Rating
// @Failure 400 {object} support.HttpError
// @Router /api/ratings [post]
func (controller ratingController) Create(c echo.Context) error {

	rating := &Rating{}
	// code := c.Param("majorcode")

	rating.Author = c.FormValue("author")
	rating.Productcode = c.FormValue("productcode")
	rating.Productname = c.FormValue("productname")
	rating.Description = c.FormValue("description")
	rating.Shopalias = c.Get("bizname").(string)
	rat := c.FormValue("rating")
	r, e := strconv.ParseInt(rat, 10, 64)
	if e != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the raitn!")
	}
	rating.Rate = r
	// fmt.Println(">>>>>>>>>>>rating create", rating)
	_, err1 := controller.service.Create(rating)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "created successifuly")
} // ///////controllers/////////////////
func (controller ratingController) Create1(c echo.Context) error {

	rating := &Rating{}
	err := c.Bind(&rating)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	// fmt.Println(">>>>>>>>>>>rating Bizname create 1", rating)
	rating, err1 := controller.service.Create(rating)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, rating)
}

// GetAll godoc
// @Summary GetAll a rating
// @Description Getall ratings
// @Tags ratings
// @Accept json
// @Produce json
// @Success 201 {object} Rating
// @Failure 400 {object} support.HttpError
// @Router /api/ratings [get]
func (controller ratingController) GetAll(c echo.Context) error {

	// fmt.Println(">>>>>>>>>>>rating Bizname", Bizname)
	ratings, err3 := controller.service.GetAll()
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, ratings)
}

// @Summary Get a rating
// @Description Get item
// @Tags ratings
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} Rating
// @Failure 400 {object} support.HttpError
// @Router /api/ratings [get]
func (controller ratingController) GetOne(c echo.Context) error {

	id := c.Param("code")
	fmt.Println(">>>>>>>>>>>rating Bizname get one=====>", c.Request().Body)
	rating, problem := controller.service.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, rating)
}
func (controller ratingController) Featured(c echo.Context) error {

	code := c.Param("code")
	status := c.FormValue("status")
	fmt.Println("ccccccccccccccccccccc", status)
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	fmt.Println("jjjjjjjjjjjjjjj", feat)
	problem := controller.service.Featured(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update a rating
// @Description Update a rating item
// @Tags ratings
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Rating
// @Failure 400 {object} support.HttpError
// @Router /api/ratings [put]
func (controller ratingController) Update(c echo.Context) error {

	rating := &Rating{}
	rating.Author = c.FormValue("author")
	rating.Productcode = c.FormValue("productcode")
	rating.Productname = c.FormValue("productname")
	rating.Description = c.FormValue("description")
	id := c.Param("code")
	// fmt.Println("----------------", id)

	r, e := strconv.ParseInt(c.FormValue("rating"), 10, 64)
	if e != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the raitn!")
	}
	rating.Rate = r

	_, problem := controller.service.Update(id, rating)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusCreated, "updated successifuly")
}

// Delete godoc
// @Summary Delete a rating
// @Description Create a new rating item
// @Tags ratings
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} string
// @Failure 400 {object} support.HttpError
// @Router /api/ratings [delete]
func (controller ratingController) Delete(c echo.Context) error {

	id := c.Param("code")
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
