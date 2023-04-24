package loan

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/myrachanto/demyst/src/api/accounting"
	"github.com/myrachanto/demyst/src/api/business"
	httperrors "github.com/myrachanto/erroring"
)

var (
	LoanService LoanServiceInterface = &loanService{}
)

type LoanServiceInterface interface {
	Create(loan *Loan) (*Loan, httperrors.HttpErr)
	GetOne(code string) (*Results, httperrors.HttpErr)
	GetAll(search string) ([]*Loan, httperrors.HttpErr)
	Submit(code string) (string, httperrors.HttpErr)
	AccountingSoftware(accounts []BalanceSheet, loan *Loan) (*Loan, httperrors.HttpErr)
	ReadJsonData(filename string) ([]BalanceSheet, httperrors.HttpErr)
	GetMonth(month string) int
	GetTheLast12MonthS(month, year int, accounts []BalanceSheet) (res []BalanceSheet)
	CheckIfGetAverageAsset12BiggerLoan(accounts []BalanceSheet, amount float64) bool
	CheckIfAllProfit12Months(accounts []BalanceSheet) bool
	AcountData(accounts []BalanceSheet, loan *Loan) (int, httperrors.HttpErr)
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
	business, err := business.Businessrepo.GetOneByName(loan.Business)
	if err != nil {
		return nil, err
	}
	loan.BusinessCode = business.Code
	loan.YearEstablished = int32(business.YearEstablished)
	accounting, err := accounting.Accountingrepo.GetOneByName(loan.AccountingSoftware)
	if err != nil {
		return nil, err
	}
	loan.AccountingSoftwareCode = accounting.Code
	return service.repo.Create(loan)
}

func (service *loanService) GetAll(search string) ([]*Loan, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}
func (service *loanService) GetOne(code string) (*Results, httperrors.HttpErr) {
	loan, err := service.repo.GetOne(code)
	if err != nil {
		return nil, err
	}
	filename := "./sampledata1.json"
	accounts, err := service.ReadJsonData(filename)
	if err != nil {
		return nil, err
	}
	loan2, err := service.AccountingSoftware(accounts, loan)
	if err != nil {
		return nil, err
	}
	return &Results{
		Data: accounts,
		Loan: *loan2,
	}, nil
}
func (service *loanService) Submit(code string) (string, httperrors.HttpErr) {
	loan, err := service.repo.GetOne(code)
	if err != nil {
		return "", err
	}
	filename := "./sampledata1.json"
	accounts, err := service.ReadJsonData(filename)
	if err != nil {
		return "", err
	}
	loan2, err := service.AccountingSoftware(accounts, loan)
	if err != nil {
		return "", err
	}
	_, err = service.repo.LoanUpdatePreassesment(loan.Code, loan2.PreAssesment)
	if err != nil {
		return "", err
	}
	res1, err := service.DecisionAlgorithym(loan2, accounts)
	if err != nil {
		return "", err
	}
	return service.repo.LoanUpdate(loan.Code, res1)
}
func (service *loanService) AccountingSoftware(accounts []BalanceSheet, loan *Loan) (*Loan, httperrors.HttpErr) {
	if !loan.ActiveSoftware {
		results, err2 := service.AcountData(accounts, loan)
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
func (service *loanService) DecisionAlgorithym(loan *Loan, accounts []BalanceSheet) (string, httperrors.HttpErr) {
	if !loan.ActiveDecision {

		if len(accounts) >= 12 {
			if loan.PreAssesment == 100 {
				return APPROVED, nil
			} else if loan.PreAssesment == 60 {
				return APPROVED, nil
			} else {
				return DECLINED, nil
			}
		} else {

			return DECLINED, nil
		}
	} else {

		// TODO write a query to algo
		algorithyendpoint := ""
		response, err := http.Get(algorithyendpoint)
		if err != nil {
			return "", httperrors.NewNotFoundError("could not reach the accounting software")
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", httperrors.NewBadRequestError("could not read the data from the accounting software")
		}
		resp := &MessageResp
		err = json.Unmarshal([]byte(body), resp)
		if err != nil {
			return "", httperrors.NewBadRequestError("could not unmarshal the data")
		}
		return resp.Message, nil
	}
}
func (service *loanService) ReadJsonData(filename string) ([]BalanceSheet, httperrors.HttpErr) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, httperrors.NewBadRequestError("could not open the file")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var accounts []BalanceSheet
	err = json.Unmarshal(byteValue, &accounts)
	if err != nil {
		return nil, httperrors.NewBadRequestError("could not unmarshal the json data")
	}
	return accounts, nil
}
func (service *loanService) AcountData(accounts []BalanceSheet, loan *Loan) (int, httperrors.HttpErr) {
	// evaluate if the average asset value across 12 months is greater than the loan amount
	ok := service.CheckIfGetAverageAsset12BiggerLoan(accounts, loan.Amount)
	if ok {
		return 100, nil
	}
	//evaluate if the business has made profit in the last 12 months
	ok = service.CheckIfAllProfit12Months(accounts)
	if ok {
		return 60, nil
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
		if account.Year == year && account.Month == month {
			results = append(results, account)
			month--
		} else {
			year--
			month = 12
			if account.Year == year && account.Month == month {
				results = append(results, account)
				month--
			}

		}

	}
	return results
} 
func (service *loanService) CheckIfGetAverageAsset12BiggerLoan(accounts []BalanceSheet, amount float64) bool {
	t := time.Now()
	month := t.Month()
	year := t.Year()
	lastMonth := service.GetMonth(month.String()) - 1
	last12Months := service.GetTheLast12MonthS(lastMonth, year, accounts)
	total := 0.0

	for _, account := range last12Months {
		// fmt.Println("vamos----------", account.AssetsValue)
		total += account.AssetsValue
	}
	results := total / 12
	// fmt.Println("--------------", total, results)
	return results > amount
}


