package loan

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/myrachanto/demyst/src/db"
	"github.com/myrachanto/demyst/src/support"
	httperrors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// loanrepository repository
var (
	Loanrepository LoanrepoInterface = &loanrepository{}
	ctx                              = context.TODO()
	loanrepo                         = loanrepository{}
)

type LoanrepoInterface interface {
	Create(loan *Loan) (*Loan, httperrors.HttpErr)
	LoanUpdate(code, status string) (string, httperrors.HttpErr)
	GetOne(id string) (*Loan, httperrors.HttpErr)
	GetAll(string) ([]*Loan, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	LoanUpdatePreassesment(code string, preasssement int) (int, httperrors.HttpErr)
}
type loanrepository struct{}

func NewloanRepo() LoanrepoInterface {
	return &loanrepository{}
}

func (r *loanrepository) Create(loan *Loan) (*Loan, httperrors.HttpErr) {
	if err1 := loan.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	loan.Code = code
	loan.Base.Updated_At = time.Now()
	loan.Base.Created_At = time.Now()
	loan.Status = PENDING
	collection := db.Mongodb.Collection("loan")
	result1, err := collection.InsertOne(ctx, &loan)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create loan Failed, %d", err))
	}
	loan.ID = result1.InsertedID.(primitive.ObjectID)
	return loan, nil
}
func (r *loanrepository) LoanUpdatePreassesment(code string, preasssement int) (int, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return 0, stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("loan")
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"preassesment", preasssement}}},
		},
	)
	if errs != nil {
		return 0, httperrors.NewNotFoundError("Error updating!")
	}
	return preasssement, nil
}
func (r *loanrepository) LoanUpdate(code, status string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("loan")
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"status", status}}},
		},
	)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return status, nil
}
func (r *loanrepository) GetOne(code string) (loan *Loan, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("loan")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&loan)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return loan, nil
}
func (r *loanrepository) GetAll(search string) ([]*Loan, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("loan")
	results := []*Loan{}
	fmt.Println(search)
	if search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search, Options: "i"}}},
				bson.D{{"business_pin", primitive.Regex{Pattern: search, Options: "i"}}},
			}},
		}
		// fmt.Println(filter)
		cursor, err := collection.Find(ctx, filter)
		fmt.Println(cursor)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		fmt.Println(results)
		return results, nil
	} else {
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		return results, nil
	}

}

func (r loanrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("loan")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "loanCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r loanrepository) getuno(code string) (result *Loan, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("loan")
	filter := bson.M{"loancode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r loanrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("loan")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
