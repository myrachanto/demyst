package loan

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var jsondatas = `{"fname":"Jane","business_pin":"1234567","year_established":2022,"amount":20000
}`
var jsondatas60 = `{"fname":"Jane","business_pin":"1234567","year_established":2022,"amount":100000
}`

func TestAccountingSoftware(t *testing.T) {
	testcases := []struct {
		name     string
		filename string
		loanjson string
		value    int
	}{
		{name: "ok", filename: "../../../sampledata1.json", loanjson: jsondatas, value: 100},
		{name: "ok60percent", filename: "../../../sampledata1.json", loanjson: jsondatas60, value: 60},
		{name: "ok20percent", filename: "../../../sampledata2.json", loanjson: jsondatas60, value: 20},
	}
	service := NewloanService(NewloanRepo())
	loan := &Loan{}

	for _, test := range testcases {
		if err := json.Unmarshal([]byte(test.loanjson), &loan); err != nil {
			t.Errorf("failed to unmarshal loan data %v", err.Error())
		}
		t.Run(test.name, func(t *testing.T) {
			accounts, err := service.ReadJsonData(test.filename)
			assert.Nil(t, err)
			pre_assesment, err := service.AccountingSoftware(accounts, loan)
			assert.Nil(t, err)
			assert.EqualValues(t, test.value, pre_assesment.PreAssesment)
		})
	}

}
func TestAcountData(t *testing.T) {
	filename := "../../../sampledata1.json"
	loan := &Loan{}
	if err := json.Unmarshal([]byte(jsondatas), &loan); err != nil {
		t.Errorf("failed to unmarshal loan data %v", err.Error())
	}
	service := NewloanService(NewloanRepo())
	accounts, err := service.ReadJsonData(filename)
	assert.Nil(t, err)
	pre_assesment, err := service.AcountData(accounts, loan)
	assert.Nil(t, err)
	assert.EqualValues(t, 100, pre_assesment)

}
func TestReadJsonData(t *testing.T) {
	filename := "../../../sampledata1.json"
	service := NewloanService(NewloanRepo())
	accounts, err := service.ReadJsonData(filename)
	assert.Nil(t, err)
	assert.NotNil(t, accounts)

}
func TestGetMonth(t *testing.T) {
	month := "October"
	service := NewloanService(NewloanRepo())
	monthNumber := service.GetMonth(month)
	assert.EqualValues(t, 10, monthNumber)

}
func TestGetTheLast12MonthS(t *testing.T) {
	service := NewloanService(NewloanRepo())
	tim := time.Now()
	month := tim.Month()
	year := tim.Year()
	lastMonth := service.GetMonth(month.String()) - 1
	filename := "../../../sampledata1.json"
	accounts, err := service.ReadJsonData(filename)
	assert.Nil(t, err)
	last12Months := service.GetTheLast12MonthS(lastMonth, year, accounts)
	assert.EqualValues(t, 12, len(last12Months))
}
func TestCheckIfGetAverageAsset12BiggerLoan(t *testing.T) {
	service := NewloanService(NewloanRepo())
	filename := "../../../sampledata1.json"
	accounts, err := service.ReadJsonData(filename)
	assert.Nil(t, err)
	ok := service.CheckIfGetAverageAsset12BiggerLoan(accounts, 20000)
	assert.EqualValues(t, true, ok)
}
func TestCheckIfAllProfit12Months(t *testing.T) {
	service := NewloanService(NewloanRepo())
	filename := "../../../sampledata1.json"
	accounts, err := service.ReadJsonData(filename)
	assert.Nil(t, err)
	ok := service.CheckIfAllProfit12Months(accounts)
	assert.EqualValues(t, true, ok)
}
