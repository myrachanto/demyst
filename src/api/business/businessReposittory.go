package business

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

// businessrepository repository
var (
	Businessrepository BusinessrepoInterface = &businessrepository{}
	ctx                                      = context.TODO()
	Businessrepo                             = businessrepository{}
)

type BusinessrepoInterface interface {
	Create(business *Business) (*Business, httperrors.HttpErr)
	GetOne(id string) (*Business, httperrors.HttpErr)
	GetAll(string) ([]*Business, httperrors.HttpErr)
	Update(code string, business *Business) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (business *Business, errors httperrors.HttpErr)
}
type businessrepository struct{}

func NewbusinessRepo() BusinessrepoInterface {
	return &businessrepository{}
}

func (r *businessrepository) Create(business *Business) (*Business, httperrors.HttpErr) {
	if err1 := business.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	business.Code = code
	business.Base.Updated_At = time.Now()
	business.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("business")
	result1, err := collection.InsertOne(ctx, &business)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create business Failed, %d", err))
	}
	business.ID = result1.InsertedID.(primitive.ObjectID)
	return business, nil
}
func (r *businessrepository) GetOne(code string) (business *Business, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("business")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&business)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return business, nil
}
func (r *businessrepository) GetOneByName(name string) (business *Business, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("business")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&business)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return business, nil
}
func (r *businessrepository) GetAll(search string) ([]*Business, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("business")
	results := []*Business{}
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

func (r *businessrepository) Update(code string, business *Business) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("business")
	var ac Business
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if business.Name == "" {
		business.Name = ac.Name
	}
	if business.Code == "" {
		business.Code = ac.Code
	}
	if business.BusinessPin == "" {
		business.BusinessPin = ac.BusinessPin
	}
	if business.YearEstablished == 0 {
		business.YearEstablished = ac.YearEstablished
	}
	update := bson.M{"$set": business}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}
func (r businessrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("business")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r businessrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("business")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "businessCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r businessrepository) getuno(code string) (result *Business, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("business")
	filter := bson.M{"businesscode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r businessrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("business")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
