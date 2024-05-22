package majorcategory

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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// majorcategoryrepository ...
var (
	Majorcategoryrepository MajorcategoryRepoInterface = &majorcategoryrepository{}
	ctx                                                = context.TODO()
	Major                                              = majorcategoryrepository{}
)

type MajorcategoryRepoInterface interface {
	Create(majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr)
	GetOne(id string) (*Majorcategory, httperrors.HttpErr)
	GetAll() ([]Majorcategory, httperrors.HttpErr)
	GetAll1(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(id string, majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	GetAllMajorcategories(Bizname string) ([]Majorcategory, httperrors.HttpErr)
	GetOnebyName(name string) (majorcategory *Majorcategory, errors httperrors.HttpErr)
}

type majorcategoryrepository struct {
	Mongodb *mongo.Database
	Bizname string
	Cancel  context.CancelFunc
}

func NewmajorcategoryRepo() MajorcategoryRepoInterface {
	return &majorcategoryrepository{}
}

func (r *majorcategoryrepository) Create(majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr) {

	if err1 := majorcategory.Validate(); err1 != nil {
		return nil, err1
	}
	code, err1 := r.genecode()
	if err1 != nil {
		return nil, err1
	}

	ok := r.CheckifExist(majorcategory.Name)
	if ok {
		return nil, httperrors.NewNotFoundError("That Major Category already exist!")
	}
	majorcategory.Base.Updated_At = time.Now()
	majorcategory.Base.Created_At = time.Now()
	// majorcategory.Url = support.Joins(majorcategory.Name)
	majorcategory.Code = code
	// fmt.Println("major cat-----------step 1", majorcategory)
	collection := db.Mongodb.Collection("majorcategory")
	_, err := collection.InsertOne(ctx, majorcategory)
	if err != nil {
		fmt.Println("major cat-----------step 2", err)
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create product Failed, %d", err))
	}
	tag, e := r.Getuno(code)

	if e != nil {
		return nil, e
	}
	return tag, nil
}

func (r *majorcategoryrepository) GetOnebyUrl(biz, name string) (majorcategory *Majorcategory, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}

	collection := db.Mongodb.Collection("majorcategory")
	filter := bson.D{
		{"$and", bson.A{
			bson.M{"url": primitive.Regex{Pattern: name, Options: "i"}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&majorcategory)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}

	return majorcategory, nil
}
func (r *majorcategoryrepository) GetOnebyName(name string) (ac *Majorcategory, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	trimmed := strings.TrimSpace(name)
	collection := db.Mongodb.Collection("majorcategory")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: trimmed, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		// fmt.Println("---------+++++++++++++", err)
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return ac, nil
}
func (r *majorcategoryrepository) GetOne(code string) (majorcategory *Majorcategory, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}

	collection := db.Mongodb.Collection("majorcategory")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&majorcategory)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}

	return majorcategory, nil
}

func (r *majorcategoryrepository) GetAll() ([]Majorcategory, httperrors.HttpErr) {
	majorcategorys := []Majorcategory{}
	collection := db.Mongodb.Collection("majorcategory")
	// fmt.Println("-------------------------getall major categories")
	filter := bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var majorcategory Majorcategory
		err := cur.Decode(&majorcategory)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		majorcategorys = append(majorcategorys, majorcategory)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	// fmt.Println("-------------------------getall major categories", majorcategorys)
	return majorcategorys, nil

}
func (r *majorcategoryrepository) GetAll1(search support.Paginator) (*Results, httperrors.HttpErr) {
	majorcategorys := []Majorcategory{}
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

	collection := db.Mongodb.Collection("majorcategory")
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
		var majorcategory Majorcategory
		err := cur.Decode(&majorcategory)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		majorcategorys = append(majorcategorys, majorcategory)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	count, _ := collection.CountDocuments(ctx, filter)
	// fmt.Println("-------------------------getall major categories", majorcategorys)
	return &Results{
		Data:        majorcategorys,
		Total:       int(count),
		Pages:       float32(count) / float32(search.Pagesize),
		CurrentPage: search.Page,
	}, nil

}

func (r *majorcategoryrepository) Update(code string, majorcategory *Majorcategory) (*Majorcategory, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	result, err3 := r.Getuno(code)
	if err3 != nil {
		fmt.Println(err3)
	}
	fmt.Println("result------------------------")
	if majorcategory.Name == "" {
		majorcategory.Name = result.Name
	}
	if majorcategory.Title == "" {
		majorcategory.Title = result.Title
	}
	if majorcategory.Description == "" {
		majorcategory.Description = result.Description
	}
	if majorcategory.Supercategory == "" {
		majorcategory.Supercategory = result.Supercategory
	}
	if majorcategory.Description == "" {
		majorcategory.Description = result.Description
	}
	if majorcategory.Code == "" {
		majorcategory.Code = result.Code
	}
	if majorcategory.Shopalias == "" {
		majorcategory.Shopalias = result.Shopalias
	}
	if majorcategory.Url == "" {
		majorcategory.Url = result.Url
	}
	if majorcategory.Meta == "" {
		majorcategory.Meta = result.Meta
	}
	if majorcategory.Content == "" {
		majorcategory.Content = result.Content
	}
	majorcategory.Base.Updated_At = time.Now()
	majorcategory.Url = support.Joins(majorcategory.Name)
	// if len(majorcategory.Category) <= 0 {
	// 	majorcategory.Category = result.Category
	// }
	collection := db.Mongodb.Collection("majorcategory")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	update := bson.M{"$set": majorcategory}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, httperrors.NewNotFoundError("Error updating!")
	}

	tag, e := r.Getuno(code)

	if e != nil {
		return nil, e
	}
	fmt.Println(tag)
	return tag, nil
}
func (r *majorcategoryrepository) Delete(code string) (string, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("majorcategory")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil
}
func (r *majorcategoryrepository) genecode() (string, httperrors.HttpErr) {

	special := support.Stamper()
	collection := db.Mongodb.Collection("majorcategory")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	code := "majorcategoryCode" + strconv.FormatUint(uint64(co), 10) + "-" + special

	code = support.Hasher(code)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r *majorcategoryrepository) Getuno(code string) (result *Majorcategory, err httperrors.HttpErr) {
	// fmt.Println("--------+++++++++", code)
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("majorcategory")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		// fmt.Println("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz>>>>>>>>>>>>>>>>>>>>")
		return nil, httperrors.NewNotFoundError("no results found")
	}

	// fmt.Println("--------+++++++++ ????????step two get 2", code)
	return result, nil
}
func (r *majorcategoryrepository) CheckifExist(name string) bool {

	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return false
	}
	result := Majorcategory{}
	collection := db.Mongodb.Collection("majorcategory")
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
func (r *majorcategoryrepository) GetAllMajorcategories(Bizname string) ([]Majorcategory, httperrors.HttpErr) {
	majorcategorys := []Majorcategory{}
	collection := db.Mongodb.Collection("majorcategory")
	// fmt.Println("-------------------------getall major categories")
	filter := bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var majorcategory Majorcategory
		err := cur.Decode(&majorcategory)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		majorcategorys = append(majorcategorys, majorcategory)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	// fmt.Println("-------------------------getall major categories", majorcategorys)
	return majorcategorys, nil

}
