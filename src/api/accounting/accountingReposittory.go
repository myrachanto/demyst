package accounting

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

// accountingrepository repository
var (
	Accountingrepository AccountingrepoInterface = &accountingrepository{}
	ctx                                          = context.TODO()
	accountingrepo                               = accountingrepository{}
)

type AccountingrepoInterface interface {
	Create(accounting *Accounting) (*Accounting, httperrors.HttpErr)
	GetOne(id string) (*Accounting, httperrors.HttpErr)
	GetAll(string) ([]*Accounting, httperrors.HttpErr)
	Update(code string, accounting *Accounting) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
}
type accountingrepository struct{}

func NewaccountingRepo() AccountingrepoInterface {
	return &accountingrepository{}
}

func (r *accountingrepository) Create(accounting *Accounting) (*Accounting, httperrors.HttpErr) {
	if err1 := accounting.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	accounting.Code = code
	accounting.Base.Updated_At = time.Now()
	accounting.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("accounting")
	result1, err := collection.InsertOne(ctx, &accounting)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create accounting Failed, %d", err))
	}
	accounting.ID = result1.InsertedID.(primitive.ObjectID)
	return accounting, nil
}
func (r *accountingrepository) GetOne(code string) (accounting *Accounting, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("accounting")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&accounting)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return accounting, nil
}
func (r *accountingrepository) GetAll(search string) ([]*Accounting, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("accounting")
	results := []*Accounting{}
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

func (r *accountingrepository) Update(code string, accounting *Accounting) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("accounting")
	var ac Accounting
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if accounting.Name == "" {
		accounting.Name = ac.Name
	}
	if accounting.Code == "" {
		accounting.Code = ac.Code
	}
	if accounting.BusinessPin == "" {
		accounting.BusinessPin = ac.BusinessPin
	}
	if accounting.UrlEndpoint == "" {
		accounting.UrlEndpoint = ac.UrlEndpoint
	}
	update := bson.M{"$set": accounting}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}
func (r accountingrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("accounting")

	filter := bson.M{"usercode": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r accountingrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("accounting")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "accountingCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r accountingrepository) getuno(code string) (result *Accounting, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("accounting")
	filter := bson.M{"accountingcode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r accountingrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("accounting")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
