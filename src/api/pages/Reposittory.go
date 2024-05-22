package pages

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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// pagerepository repository
var (
	Pagerepository PagerepoInterface = &pagerepository{}
	ctx                              = context.TODO()
	Pagerepo                         = pagerepository{}
)

type PagerepoInterface interface {
	Create(page *Page) (*Page, httperrors.HttpErr)
	GetOne(id string) (*Page, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, page *Page) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *Page, errors httperrors.HttpErr)
	GetOneByUrl(code string) (page *Page, errors httperrors.HttpErr)
}
type pagerepository struct{}

func NewpageRepo() PagerepoInterface {
	return &pagerepository{}
}

func (r *pagerepository) Create(page *Page) (*Page, httperrors.HttpErr) {
	if err1 := page.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	page.Code = code
	page.Base.Updated_At = time.Now()
	page.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("page")
	result1, err := collection.InsertOne(ctx, &page)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create page Failed, %d", err))
	}
	page.ID = result1.InsertedID.(primitive.ObjectID)
	return page, nil
}
func (r *pagerepository) GetOne(code string) (page *Page, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("page")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return page, nil
}

func (r *pagerepository) GetOneByUrl(code string) (page *Page, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("page")
	filter := bson.M{"name": code}
	err := collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return page, nil
}
func (r *pagerepository) GetOneByUrl2(code string) (page *Page, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("page")
	filter := bson.M{"url": code}
	err := collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return page, nil
}
func (r *pagerepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.D{{"name", -1}})
	collection := db.Mongodb.Collection("page")
	results := []*Page{}
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

func (r *pagerepository) Update(code string, page *Page) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("page")
	var ac Page
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if page.Name == "" {
		page.Name = ac.Name
	}
	if page.Code == "" {
		page.Code = ac.Code
	}
	if page.Title == "" {
		page.Title = ac.Title
	}
	if page.Meta == "" {
		page.Meta = ac.Meta
	}
	if page.Url == "" {
		page.Url = ac.Url
	}
	if page.Content == "" {
		page.Content = ac.Content
	}
	update := bson.M{"$set": page}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r *pagerepository) GetOneByName(name string) (ac *Page, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("page")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return ac, nil
}
func (r pagerepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("page")
	// fmt.Println("=====================asacaa deleting page", id)
	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r pagerepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("page")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "pageCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r pagerepository) getuno(code string) (result *Page, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("page")
	filter := bson.M{"pagecode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r pagerepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("page")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
