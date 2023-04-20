package loan

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	httperrors "github.com/myrachanto/erroring"
)

var (
	LoanService LoanServiceInterface = &loanService{}
)

type LoanServiceInterface interface {
	Create(loan *Loan) (*Loan, httperrors.HttpErr)
	GetOne(code string) (loan *Loan, errors httperrors.HttpErr)
	GetAll(search string) ([]*Loan, httperrors.HttpErr)
}
type loanService struct {
	repo LoanrepoInterface
}

func NewloanService(repository LoanrepoInterface) LoanServiceInterface {
	return &loanService{
		repository,
	}
}
func (service *loanService) Create(loan *Loan) (*Loan, httperrors.HttpErr) {

	loan1, err := service.repo.Create(loan)
	loan2, err := service.AccountingSoftware(loan1)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (service *loanService) GetAll(search string) ([]*Loan, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *loanService) GetOne(code string) (*Loan, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *loanService) AccountingSoftware(loan *Loan) (*Loan, httperrors.HttpErr) {
	if loan.AccountingSoftware == "" {
		resp := []BalanceSheet{}
		//TODO populate the balance sheet with the sample data
		results, err2 := service.AcountData(resp, loan)
		if err2 != nil {
			return nil, err2
		}
		loan.PreAssesment = results
		return loan, nil
	} else {
		//TODO create accounting softwares
		response, err := http.Get(loan.AccountingSoftware)
		if err != nil {
			return nil, httperrors.NewNotFoundError("could not reach the accounting software")
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, httperrors.NewBadRequestError("could not read the data from the accounting software")
		}
		var resp []BalanceSheet
		err = json.Unmarshal([]byte(body), &resp)
		if err != nil {
			return nil, httperrors.NewBadRequestError("could not unmarshal the data")
		}
		results, err2 := service.AcountData(resp, loan)
		if err2 != nil {
			return nil, err2
		}
		loan.PreAssesment = results
		return loan, nil
	}
}
func (service *loanService) DecisionAlgorithym(loan *Loan) (int, httperrors.HttpErr) {

}
func (service *loanService) AcountData(accounts []BalanceSheet, loan *Loan) (int, httperrors.HttpErr) {
	//evaluate if the business has made profit in the last 12 months
	ok := service.CheckIfAllProfit12Months(accounts)
	if ok {
		return 60, nil
	}
	// evaluate if the average asset value across 12 months is greater than the loan amount
	ok = service.CheckIfGetAverageAsset12(accounts, loan.Amount)
	if ok {
		return 100, nil
	}
	return 20, nil
}
func (service *loanService) CheckIfAllProfit12Months(accounts []BalanceSheet) bool {
	var res bool = true
	t := time.Now()
	month := t.Month()
	year := t.Year()
	lastMonth := service.GetMonth(month.String()) - 1
	last12Months := service.GetTheLast12MonthS(lastMonth, year, accounts)
	for _, account := range last12Months {
		if account.ProfitOrLoss < 0 {
			res = false
			break
		}
	}

	return res
}
func (service *loanService) CheckIfGetAverageAsset12(accounts []BalanceSheet, amount float64) bool {
	t := time.Now()
	month := t.Month()
	year := t.Year()
	lastMonth := service.GetMonth(month.String()) - 1
	last12Months := service.GetTheLast12MonthS(lastMonth, year, accounts)
	total := 0.0
	for _, account := range last12Months {
		total += account.AssetsValue
	}
	results := total / 12

	return results > amount
}

func (service *loanService) GetMonth(month string) int {
	monthslookup := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	var res = -1
	for k, v := range monthslookup {
		if month == v {
			res = k + 1
			break
		}
	}
	return res
}
func (service *loanService) GetTheLast12MonthS(month, year int, accounts []BalanceSheet) (res []BalanceSheet) {
	results := []BalanceSheet{}
	for _, account := range accounts {
		if len(results) == 12 {
			break
		}
		if month > 0 {
			if account.Year == year && account.Month == month {
				results = append(results, account)
			}
		} else {
			year = year - 1
			if account.Year == year && account.Month == month {
				results = append(results, account)
			}
		}

	}
	return results
}
