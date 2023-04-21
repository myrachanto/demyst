package loan

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// loanController ...
var (
	LoanController loanControllerInterface = &loanController{}
)

type loanControllerInterface interface {
	Create(c echo.Context) error
	GetOne(c echo.Context) error
	GetAll(c echo.Context) error
	Submit(c echo.Context) error
}

type loanController struct {
	service LoanServiceInterface
}

func NewloanController(ser LoanServiceInterface) loanControllerInterface {
	return &loanController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a loan
// @Description Create a new loan item
// @Tags loans
// @Accept json
// @Produce json
// @Success 201 {object} loan
// @Failure 400 {object} support.HttpError
// @Router /api/loans [post]
func (controller loanController) Create(c echo.Context) error {

	loan := &Loan{}
	loan.Name = c.FormValue("name")
	loan.BusinessPin = c.FormValue("business_pin")
	loan.AccountingSoftware = c.FormValue("accounting_software")
	amount := c.FormValue("amount")
	loanAmount, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		return c.JSON(http.StatusCreated, "failed to parse the loan amount")
	}
	loan.Amount = loanAmount
	yoe := c.FormValue("year_established")
	yearofEstablishment, err := strconv.ParseUint(yoe, 10, 32)
	if err != nil {
		return c.JSON(http.StatusCreated, "failed to parse the loan amount")
	}

	loan.YearEstablished = int32(yearofEstablishment)

	_, err1 := controller.service.Create(loan)
	if err1 != nil {
		return c.JSON(err1.Code(), err1.Message())
	}
	return c.JSON(http.StatusCreated, "loan created succesifully")

}

// GetAll godoc
// @Summary GetAll a loan
// @Description Getall loans
// @Tags loans
// @Accept json
// @Produce json
// @Success 201 {object} loan
// @Failure 400 {object} support.HttpError
// @Router /api/loans [get]
func (controller loanController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")
	loans, err3 := controller.service.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, loans)
}

// @Summary Get a loan
// @Description Get item
// @Tags loans
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} loan
// @Failure 400 {object} support.HttpError
// @Router /api/loans [get]
func (controller loanController) GetOne(c echo.Context) error {
	code := c.Param("code")
	loan, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, loan)
}

// @Summary Submit a loan
// @Description Submit item
// @Tags loans
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} loan
// @Failure 400 {object} support.HttpError
// @Router /api/loans [get]
func (controller loanController) Submit(c echo.Context) error {
	code := c.Param("code")
	loan, problem := controller.service.Submit(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, loan)
}
