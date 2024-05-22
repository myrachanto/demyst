package feature

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

// featurerepository repository
var (
	Featurerepository FeaturerepoInterface = &featurerepository{}
	ctx                                    = context.TODO()
	Featurerepo                            = featurerepository{}
)

type FeaturerepoInterface interface {
	Create(feature *Feature) (*Feature, httperrors.HttpErr)
	GetOne(id string) (*Feature, httperrors.HttpErr)
	GetAll(string) ([]Feature, httperrors.HttpErr)
	GetAll1(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, feature *Feature) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	// Count() (float64, httperrors.HttpErr)
	// GetOneByName(name string) (ac *Feature, errors httperrors.HttpErr)
}
type featurerepository struct{}

func NewfeatureRepo() FeaturerepoInterface {
	return &featurerepository{}
}

func (r *featurerepository) Create(feature *Feature) (*Feature, httperrors.HttpErr) {
	if err1 := feature.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	feature.Code = code
	feature.Base.Updated_At = time.Now()
	feature.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("feature")
	result1, err := collection.InsertOne(ctx, &feature)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create feature Failed, %d", err))
	}
	feature.ID = result1.InsertedID.(primitive.ObjectID)
	return feature, nil
}
func (r *featurerepository) GetOne(code string) (feature *Feature, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("feature")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&feature)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return feature, nil
}
func (r *featurerepository) GetAll(search string) ([]Feature, httperrors.HttpErr) {
	// skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(10)
	collection := db.Mongodb.Collection("feature")
	results := []Feature{}
	fmt.Println(search)
	if search != "" {
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search, Options: "i"}}},
			}},
		}
		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		return results, nil
	} else {
		cursor, err := collection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}

		return results, nil
	}

}
func (r *featurerepository) GetAll1(search support.Paginator) (*Results, httperrors.HttpErr) {
	features := []Feature{}
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

	collection := db.Mongodb.Collection("feature")
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
		var feature Feature
		err := cur.Decode(&feature)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		features = append(features, feature)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	count, _ := collection.CountDocuments(ctx, filter)
	// fmt.Println("-------------------------getall major categories", majorfeatures)
	return &Results{
		Data:        features,
		Total:       int(count),
		Pages:       int(count) / int(search.Pagesize),
		CurrentPage: search.Page,
	}, nil

}

func (r *featurerepository) Update(code string, feature *Feature) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("feature")
	var ac Feature
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if feature.Name == "" {
		feature.Name = ac.Name
	}
	if feature.Code == "" {
		feature.Code = ac.Code
	}
	if feature.Title == "" {
		feature.Title = ac.Title
	}
	if feature.Description == "" {
		feature.Description = ac.Description
	}
	update := bson.M{"$set": feature}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r *featurerepository) GetOneByName(name string) (ac *Feature, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	fmt.Println("===============++++++", name)
	trimmed := strings.TrimSpace(name)
	collection := db.Mongodb.Collection("feature")

	fmt.Println("===============++++++", name)
	// filter := bson.D{
	// 	{"$and", bson.A{
	// 		bson.D{{"name", trimmed}},
	// 	}},
	// }
	filter := bson.M{"name": trimmed}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find feature with this id, %s", name))
	}
	return ac, nil
}

func (r featurerepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("feature")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r featurerepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("feature")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "featureCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r featurerepository) getuno(code string) (result *Feature, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("feature")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}

func (r featurerepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("feature")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
