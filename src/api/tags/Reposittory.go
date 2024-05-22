package tags

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

// tagrepository repository
var (
	Tagrepository TagrepoInterface = &tagrepository{}
	ctx                            = context.TODO()
	Tagrepo                        = tagrepository{}
)

type TagrepoInterface interface {
	Create(tag *Tag) (*Tag, httperrors.HttpErr)
	GetOne(id string) (*Tag, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	GetAll1(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, tag *Tag) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *Tag, errors httperrors.HttpErr)
	GetOneByUrl(name string) (ac *Tag, errors httperrors.HttpErr)
}
type tagrepository struct{}

func NewtagRepo() TagrepoInterface {
	return &tagrepository{}
}

func (r *tagrepository) Create(tag *Tag) (*Tag, httperrors.HttpErr) {
	if err1 := tag.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	tag.Code = code
	tag.Base.Updated_At = time.Now()
	tag.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("tag")
	result1, err := collection.InsertOne(ctx, &tag)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create tag Failed, %d", err))
	}
	tag.ID = result1.InsertedID.(primitive.ObjectID)
	return tag, nil
}
func (r *tagrepository) GetOne(code string) (tag *Tag, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("tag")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&tag)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return tag, nil
}
func (r *tagrepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.D{{"name", -1}})
	collection := db.Mongodb.Collection("tag")
	results := []Tag{}
	fmt.Println(search)
	if search.Search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search.Search, Options: "i"}}},
				bson.D{{"business_pin", primitive.Regex{Pattern: search.Search, Options: "i"}}},
			}},
		}
		// fmt.Println(filter)
		cursor, err := collection.Find(ctx, filter, findOptions)
		fmt.Println(cursor)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		return &Results{
			Data:  results,
			Total: len(results),
		}, nil
	} else {
		cursor, err := collection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}

		count, errd := r.Count()
		if errd != nil {
			return nil, errd
		}
		return &Results{
			Data:  results,
			Total: int(count),
		}, nil
	}

}

func (r *tagrepository) GetAll1(search support.Paginator) (*Results, httperrors.HttpErr) {
	tags := []Tag{}
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

	collection := db.Mongodb.Collection("tag")
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
		var tag Tag
		err := cur.Decode(&tag)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		tags = append(tags, tag)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	count, _ := collection.CountDocuments(ctx, filter)
	// fmt.Println("-------------------------getall major categories", majorcategorys)
	return &Results{
		Data:        tags,
		Total:       int(count),
		Pages:       int(count) / int(search.Pagesize),
		CurrentPage: search.Page,
	}, nil

}

func (r *tagrepository) Update(code string, tag *Tag) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("tag")
	var ac Tag
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if tag.Name == "" {
		tag.Name = ac.Name
	}
	if tag.Code == "" {
		tag.Code = ac.Code
	}
	if tag.Title == "" {
		tag.Title = ac.Title
	}
	if tag.Description == "" {
		tag.Description = ac.Description
	}
	update := bson.M{"$set": tag}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r *tagrepository) GetOneByUrl(name string) (ac *Tag, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("tag")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"url", primitive.Regex{Pattern: name, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return ac, nil
}
func (r *tagrepository) GetOneByName(name string) (ac *Tag, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	trimmed := strings.TrimSpace(name)
	collection := db.Mongodb.Collection("tag")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", trimmed}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find tag with this id, %d", err))
	}
	return ac, nil
}
func (r tagrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("tag")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r tagrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("tag")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "tagCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r tagrepository) getuno(code string) (result *Tag, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("tag")
	filter := bson.M{"tagcode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r tagrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("tag")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
