package subLocation

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// subLocationrepository repository
var (
	SubLocationrepository SubLocationrepoInterface = &subLocationrepository{}
	ctx                                            = context.TODO()
	SubLocationrepo                                = subLocationrepository{}
)

type SubLocationrepoInterface interface {
	Create(subLocation *SubLocation) (*SubLocation, httperrors.HttpErr)
	GetOne(id string) (*SubLocation, httperrors.HttpErr)
	GetAll() ([]SubLocation, httperrors.HttpErr)
	GetAll1(search support.Paginator) (*Results, httperrors.HttpErr)

	GetAllByLocation(code string) ([]SubLocation, httperrors.HttpErr)
	Update(code string, subLocation *SubLocation) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *SubLocation, errors httperrors.HttpErr)
}
type subLocationrepository struct{}

func NewsubLocationRepo() SubLocationrepoInterface {
	return &subLocationrepository{}
}

func (r *subLocationrepository) Create(subLocation *SubLocation) (*SubLocation, httperrors.HttpErr) {
	fmt.Println("step 1 --------------")
	if err1 := subLocation.Validate(); err1 != nil {
		return nil, err1
	}
	fmt.Println("step 2 --------------")
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	fmt.Println("step 3 --------------", subLocation.Location)
	location, errs := location.Locationrepo.Getuno(subLocation.Location)
	if errs != nil {
		return nil, errs
	}
	fmt.Println("step 4 --------------")
	subLocation.Location = location.Code
	subLocation.Code = code
	subLocation.Base.Updated_At = time.Now()
	subLocation.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("subLocation")
	result1, err := collection.InsertOne(ctx, &subLocation)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create subLocation Failed, %d", err))
	}
	subLocation.ID = result1.InsertedID.(primitive.ObjectID)
	return subLocation, nil
}
func (r *subLocationrepository) GetOne(code string) (subLocation *SubLocation, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("subLocation")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&subLocation)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}

	location, errs := location.Locationrepo.Getuno(subLocation.Location)
	if errs != nil {
		return nil, errs
	}
	subLocation.Location = location.Name
	return subLocation, nil
}
func (r *subLocationrepository) GetAllByLocation(code string) ([]SubLocation, httperrors.HttpErr) {
	// fmt.Println("++++++++++++++++++++++++++++++", code)
	results := []SubLocation{}
	collection := db.Mongodb.Collection("subLocation")
	var filter primitive.M
	if code == "" {
		// fmt.Println("++++++++++++++++++++++++++++++dfggg", code)
		filter = bson.M{}
	} else {
		filter = bson.M{"location": code}
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	locas := []SubLocation{}
	for _, r := range results {
		location, errs := location.Locationrepo.Getuno(r.Location)
		if errs != nil {
			return nil, errs
		}
		r.Location = location.Name
		locas = append(locas, r)
	}
	// fmt.Println("-------------------sublocations", locas)
	return locas, nil
}
func (r *subLocationrepository) GetAll() ([]SubLocation, httperrors.HttpErr) {

	results := []SubLocation{}
	collection := db.Mongodb.Collection("subLocation")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}

	locas := []SubLocation{}
	for _, r := range results {
		location, errs := location.Locationrepo.Getuno(r.Location)
		if errs != nil {
			return nil, errs
		}
		r.Location = location.Name
		locas = append(locas, r)
	}
	return locas, nil

}
func (r *subLocationrepository) GetAll1(search support.Paginator) (*Results, httperrors.HttpErr) {
	subLocations := []SubLocation{}
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

	collection := db.Mongodb.Collection("subLocation")
	// fmt.Println("-------------------------getall major categories")
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: search.Search, Options: "i"}}},
			bson.D{{"title", primitive.Regex{Pattern: search.Search, Options: "i"}}},
		}},
	}
	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var subLocation SubLocation
		err := cur.Decode(&subLocation)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		subLocations = append(subLocations, subLocation)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	count, _ := collection.CountDocuments(ctx, filter)
	// fmt.Println("-------------------------getall major categories", len(subLocations))
	locas := []SubLocation{}
	for _, r := range subLocations {
		location, errs := location.Locationrepo.Getuno(r.Location)
		if errs != nil {
			return nil, errs
		}
		r.Location = location.Name
		locas = append(locas, r)
	}
	// fmt.Println("-------------------------getall major categories", len(subLocations))
	return &Results{
		Data:        locas,
		Total:       int(count),
		Pages:       int(count) / int(search.Pagesize),
		CurrentPage: search.Page,
	}, nil

}

func (r *subLocationrepository) Update(code string, subLocation *SubLocation) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("subLocation")
	var ac SubLocation
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if subLocation.Name == "" {
		subLocation.Name = ac.Name
	}
	if subLocation.Code == "" {
		subLocation.Code = ac.Code
	}
	if subLocation.Title == "" {
		subLocation.Title = ac.Title
	}
	if subLocation.Description == "" {
		subLocation.Description = ac.Description
	}
	if subLocation.Location == "" {
		subLocation.Location = ac.Location
	}
	if subLocation.Meta == "" {
		subLocation.Meta = ac.Meta
	}
	if subLocation.Content == "" {
		subLocation.Content = ac.Content
	}
	if subLocation.PropertyType == "" {
		subLocation.PropertyType = ac.PropertyType
	}
	update := bson.M{"$set": subLocation}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r *subLocationrepository) GetOneByName(name string) (ac *SubLocation, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	trimmed := strings.TrimSpace(name)
	collection := db.Mongodb.Collection("subLocation")
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"name", trimmed}},
			bson.D{{"url", trimmed}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return ac, nil
}

func (r subLocationrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("subLocation")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r subLocationrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("subLocation")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "subLocationCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r subLocationrepository) Getuno(code string) (result *SubLocation, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("subLocation")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}

func (r subLocationrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("subLocation")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
