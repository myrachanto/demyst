package seo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/api/pages"
	subLocation "github.com/myrachanto/estate/src/api/sublocation"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// seorepository repository
var (
	Seorepository SeorepoInterface = &seorepository{}
	ctx                            = context.TODO()
	Seorepo                        = seorepository{}
)

type SeorepoInterface interface {
	Create(seo *Seo) (*Seo, httperrors.HttpErr)
	GetOne(id string) (*Seo, httperrors.HttpErr)
	GetAll(string) ([]Seo, httperrors.HttpErr)
	GetAll1(search support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, seo *Seo) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	// Count() (float64, httperrors.HttpErr)
	// GetOneByName(name string) (ac *Seo, errors httperrors.HttpErr)
}
type seorepository struct{}

func NewseoRepo() SeorepoInterface {
	return &seorepository{}
}

func (r *seorepository) Create(seo *Seo) (*Seo, httperrors.HttpErr) {
	if err1 := seo.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	seo.Code = code
	if seo.Location != "" {
		locat, errs := location.Locationrepo.GetOneByName(seo.Location)
		if errs != nil {
			return nil, errs
		}
		seo.Location = locat.Code
	}
	if seo.Sublocation != "" {
		subloc, errs := subLocation.SubLocationrepo.GetOneByName(seo.Sublocation)
		if errs != nil {
			return nil, errs
		}
		seo.Sublocation = subloc.Code
	}
	page, errd := pages.Pagerepo.GetOneByUrl2(seo.Page)
	if errd != nil {
		return nil, errd
	}
	seo.Page = page.Url
	seo.Base.Updated_At = time.Now()
	seo.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("seo")
	result1, err := collection.InsertOne(ctx, &seo)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create seo Failed, %d", err))
	}
	seo.ID = result1.InsertedID.(primitive.ObjectID)
	return seo, nil
}
func (r *seorepository) GetOne(code string) (seo *Seo, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("seo")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&seo)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return r.Assist(seo), nil
}

func (r *seorepository) Assist(seo *Seo) *Seo {
	// locaName := ""
	locat, _ := location.Locationrepo.Getuno(seo.Location)
	if locat != nil {
		seo.Location = locat.Name
	}
	sublocat, _ := subLocation.SubLocationrepo.Getuno(seo.Sublocation)
	if sublocat != nil {
		seo.Sublocation = sublocat.Name
	}
	return seo
}
func (r *seorepository) GetAll(search string) ([]Seo, httperrors.HttpErr) {
	// skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(10)
	collection := db.Mongodb.Collection("seo")
	results := []Seo{}
	var filter primitive.D
	if search != "" {
		filter = bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search, Options: "i"}}},
			}},
		}
	} else {
		filter = bson.D{}
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	finalResults := []Seo{}
	for _, res := range results {
		resp := r.Assist(&res)
		finalResults = append(finalResults, *resp)
	}
	return finalResults, nil
}
func (r *seorepository) GetAll1(search support.Paginator) (*Results, httperrors.HttpErr) {
	seos := []Seo{}
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

	collection := db.Mongodb.Collection("seo")
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
		var seo Seo
		err := cur.Decode(&seo)
		if err != nil {
			return nil, httperrors.NewNotFoundError("Error while decoding results!")
		}
		seos = append(seos, seo)
	}
	if err := cur.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Error with cursor!")
	}

	count, _ := collection.CountDocuments(ctx, filter)

	finalResults := []Seo{}
	for _, res := range seos {
		resp := r.Assist(&res)
		finalResults = append(finalResults, *resp)
	}
	// return finalResults, nil
	// fmt.Println("-------------------------getall major categories", majorseos)
	return &Results{
		Data:        finalResults,
		Total:       int(count),
		Pages:       int(count) / int(search.Pagesize),
		CurrentPage: search.Page,
	}, nil

}

func (r *seorepository) Update(code string, seo *Seo) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	// fmt.Println("step 1---------")
	if seo.Location != "" {
		locat, errs := location.Locationrepo.GetOneByName(seo.Location)
		if errs != nil {
			return "", errs
		}
		seo.Location = locat.Code
	}
	// fmt.Println("step 2---------")
	if seo.Sublocation != "" {
		subloc, errs := subLocation.SubLocationrepo.GetOneByName(seo.Sublocation)
		if errs != nil {
			return "", errs
		}
		seo.Sublocation = subloc.Code
	}
	// fmt.Println("step 3---------")
	// if seo.Page != "" {
	// 	page, errd := pages.Pagerepo.GetOne(seo.Page)
	// 	if errd != nil {
	// 		return "", errd
	// 	}
	// 	seo.Page = page.Url
	// }
	// fmt.Println("step 4---------", seo)
	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("seo")
	var ac Seo
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if seo.Name == "" {
		seo.Name = ac.Name
	}
	if seo.Code == "" {
		seo.Code = ac.Code
	}
	if seo.Title == "" {
		seo.Title = ac.Title
	}
	if seo.Description == "" {
		seo.Description = ac.Description
	}
	if seo.Kind == "" {
		seo.Kind = ac.Kind
	}
	if seo.Location == "" {
		seo.Location = ac.Location
	}
	if seo.Sublocation == "" {
		seo.Sublocation = ac.Sublocation
	}
	if seo.Page == "" {
		seo.Page = ac.Page
	}
	update := bson.M{"$set": seo}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r *seorepository) GetOneByName(name string, options ...string) *Seo {
	var results Seo
	collection := db.Mongodb.Collection("seo")
	// if
	// name = strings.Join(strings.Split(name, "-"), " ")
	fmt.Println("---------------------vamoa")
	var filter primitive.D
	if len(options) == 1 {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"page", name}},
				bson.D{{"kind", options[0]}},
			}},
		}
	} else if len(options) == 2 {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"page", name}},
				bson.D{{"kind", options[0]}},
				bson.D{{"location", options[1]}},
			}},
		}
	} else {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"page", name}},
			}},
		}
	}
	fmt.Println("d-------------------", filter)
	// filter := bson.M{"page": name}
	err := collection.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		results.Title = "eloi developers"
		results.Description = "eloi developers"
		results.Meta = "eloi developers"
	}
	// fmt.Println("d-------------------", results)
	return r.Assist(&results)
}

func (r seorepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("seo")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r seorepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("seo")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "seoCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r seorepository) getuno(code string) (result *Seo, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("seo")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return r.Assist(result), nil
}

func (r seorepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("seo")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
