package location

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// locationrepository repository
var (
	Locationrepository locationrepoInterface = &locationrepository{}
	ctx                                      = context.TODO()
	Locationrepo                             = locationrepository{}
)

type locationrepoInterface interface {
	Create(location *Location) (*Location, httperrors.HttpErr)
	GetOne(id string) (*Location, httperrors.HttpErr)
	GetAll() ([]Location, httperrors.HttpErr)
	GetAll1(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, location *Location) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *Location, errors httperrors.HttpErr)
}
type locationrepository struct{}

func NewlocationRepo() locationrepoInterface {
	return &locationrepository{}
}

func (r *locationrepository) Create(location *Location) (*Location, httperrors.HttpErr) {
	if err1 := location.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	location.Code = code
	location.Base.Updated_At = time.Now()
	location.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("location")
	result1, err := collection.InsertOne(ctx, &location)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create location Failed, %d", err))
	}
	location.ID = result1.InsertedID.(primitive.ObjectID)
	return location, nil
}
func (r *locationrepository) GetOne(code string) (location *Location, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("location")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return location, nil
}
func (r *locationrepository) GetAll() ([]Location, httperrors.HttpErr) {
	// skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(100)
	collection := db.Mongodb.Collection("location")
	results := []Location{}
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}

	return results, nil

}
func (r *locationrepository) GetAll1(search support.Paginator) (*Results, httperrors.HttpErr) {
	locations := []Location{}
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions.SetSkip(int64(skipNum))
	if search.Orderby == "name" {
		findOptions.SetSort(bson.D{{"name", search.Order}})
	}
	if search.Orderby == "title" {
		findOptions.SetSort(bson.D{{"title", search.Order}})
	}
	if search.Orderby == "" {
		findOptions.SetSort(bson.D{{"name", 1}})
	}

	collection := db.Mongodb.Collection("location")
	// fmt.Println("-------------------------getall major categories")
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: search.Search, Options: "i"}}},
			bson.D{{"title", primitive.Regex{Pattern: search.Search, Options: "i"}}},
			bson.D{{"description", primitive.Regex{Pattern: search.Search, Options: "i"}}},
		}},
	}
	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var location Location
		err := cur.Decode(&location)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		locations = append(locations, location)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	count, _ := collection.CountDocuments(ctx, filter)
	// fmt.Println("-------------------------getall major categories", majorlocations)
	return &Results{
		Data:        locations,
		Total:       int(count),
		Pages:       int(count) / int(search.Pagesize),
		CurrentPage: search.Page,
	}, nil

}

func (r *locationrepository) Update(code string, location *Location) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("location")
	var ac Location
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if location.Name == "" {
		location.Name = ac.Name
	}
	if location.Code == "" {
		location.Code = ac.Code
	}
	if location.Title == "" {
		location.Title = ac.Title
	}
	if location.Description == "" {
		location.Description = ac.Description
	}
	if location.Picture == "" {
		location.Picture = ac.Picture
	}
	if location.Meta == "" {
		location.Meta = ac.Meta
	}
	if location.Content == "" {
		location.Content = ac.Content
	}
	if location.PropertyType == "" {
		location.PropertyType = ac.PropertyType
	}
	update := bson.M{"$set": location}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r *locationrepository) GetOneByName(name string) (ac *Location, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	trimmed := strings.TrimSpace(name)
	collection := db.Mongodb.Collection("location")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: trimmed, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return ac, nil
}

func (r locationrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("location")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r locationrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("location")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "locationCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r locationrepository) Getuno(code string) (result *Location, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("location")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}

func (r locationrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("location")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
