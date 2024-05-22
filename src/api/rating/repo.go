package rating

import (
	"context"
	"fmt"
	"strconv"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ratingrepository ...
var (
	Ratingrepository RatingRepoInterface = &ratingrepository{}
	ctx                                  = context.TODO()
	Ratos                                = ratingrepository{}
)

type RatingRepoInterface interface {
	Create(rating *Rating) (*Rating, httperrors.HttpErr)
	GetOne(id string) (*Rating, httperrors.HttpErr)
	GetAll() ([]*Rating, httperrors.HttpErr)
	Featured(code string, status bool) httperrors.HttpErr
	Feature() ([]Rating, httperrors.HttpErr)
	Update(id string, rating *Rating) (*Rating, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	GetunobyName(name string) (result *Rating, err httperrors.HttpErr)
	GetRatings(code string) (result []Rating, err httperrors.HttpErr)
}

type ratingrepository struct {
	Mongodb *mongo.Database
	Cancel  context.CancelFunc
}

func NewratingRepo() RatingRepoInterface {
	return &ratingrepository{}
}

func (r *ratingrepository) Create(rating *Rating) (*Rating, httperrors.HttpErr) {
	if err1 := rating.Validate(); err1 != nil {
		return nil, err1
	}

	// fmt.Println("===================>>>>>>>", rating)
	code, err1 := r.genecode()
	if err1 != nil {
		return nil, err1
	}

	rating.Base.Updated_At = time.Now()
	rating.Base.Created_At = time.Now()
	rating.Code = code
	// fmt.Println("===================>>>>>>>", rating)
	collection := db.Mongodb.Collection("rating")
	result, err := collection.InsertOne(ctx, rating)
	if err != nil {
		fmt.Println("err -------", err)
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create product Failed, %d", err))
	}
	rating.ID = result.InsertedID.(primitive.ObjectID)
	return rating, nil
}

func (r *ratingrepository) GetOne(id string) (rating *Rating, errors httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", id}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&rating)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return rating, nil
}

func (r *ratingrepository) GetAll() ([]*Rating, httperrors.HttpErr) {
	ratings := []*Rating{}
	// fmt.Println("rating get all-------------", db.Mongodb)
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var rating Rating
		err := cur.Decode(&rating)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		ratings = append(ratings, &rating)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}
	return ratings, nil

}
func (r *ratingrepository) Feature() ([]Rating, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("rating")
	// filter := bson.M{"featured": true}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"featured", true}},
		}},
	}
	ratings := []Rating{}
	options := options.Find()
	// options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// fmt.Println(products)
	return ratings, nil

}
func (r *ratingrepository) Featured(code string, status bool) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Rating{}

	collection := db.Mongodb.Collection("rating")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"featured", status}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}

func (r *ratingrepository) Update(code string, rating *Rating) (*Rating, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	// fmt.Println("result+++++++++++++++++++++++++++")
	result, err3 := r.GetunobyCode(code, )
	if err3 != nil {
		fmt.Println(err3)
	}
	// fmt.Println("result+++++++++++++++++++++++++++step1", result)
	if rating.Productcode == "" {
		rating.Productcode = result.Productcode
	}
	if rating.Productcode == "" {
		rating.Productcode = result.Productcode
	}
	if rating.Description == "" {
		rating.Description = result.Description
	}
	if rating.Author == "" {
		rating.Author = result.Author
	}
	if rating.Code == "" {
		rating.Code = result.Code
	}
	if rating.Shopalias == "" {
		rating.Shopalias = result.Shopalias
	}
	if !rating.Featured {
		rating.Featured = result.Featured
	}
	rating.Base.Updated_At = time.Now()
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	// fmt.Println("result+++++++++++++++++++++++++++step2")
	update := bson.M{"$set": rating}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, httperrors.NewNotFoundError("Error updating!")
	}

	return rating, nil
}
func (r *ratingrepository) Delete(id string) (string, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", id}},
		}},
	}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil
}
func (r *ratingrepository) genecode() (string, httperrors.HttpErr) {

	special := support.Stamper()
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		code := "ratingCode"  + strconv.FormatUint(uint64(co), 10) + "-" + special
		return code, nil
	}
	code := "ratingCode"  + strconv.FormatUint(uint64(co), 10) + "-" + special
	return code, nil
}
func (r *ratingrepository) getuno(code string) (result *Rating, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r *ratingrepository) GetRatings(code string) (result []Rating, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	ratings := []Rating{}
	// fmt.Println("rating get all-------------", db.Mongodb)
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"productcode", code}},
		}},
	}
	cur, errs := collection.Find(ctx, filter)
	if errs != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var rating Rating
		err := cur.Decode(&rating)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		ratings = append(ratings, rating)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}
	return ratings, nil
}

func (r *ratingrepository) GetunobyCode(name string) (result *Rating, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", name}},
		}},
	}
	// fmt.Println("---------------------------", filter)
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}

func (r *ratingrepository) GetunobyName(name string) (result *Rating, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", name}},
		}},
	}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r *ratingrepository) CheckifExist(name string) bool {

	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return false
	}
	result := Rating{}
	collection := db.Mongodb.Collection("rating")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", name}},
		}},
	}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return false
	}

	return true
}
