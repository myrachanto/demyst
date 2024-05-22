package gift

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

// Giftrepository repository
var (
	Giftrepository GiftrepoInterface = &giftrepository{}
	ctx                              = context.TODO()
	Giftrepo                         = giftrepository{}
)

type GiftrepoInterface interface {
	Create(Gift *Gift) (*Gift, httperrors.HttpErr)
	GetOne(id string) (*Gift, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, Gift *Gift) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *Gift, errors httperrors.HttpErr)
}
type giftrepository struct{}

func NewGiftRepo() GiftrepoInterface {
	return &giftrepository{}
}

func (r giftrepository) Create(Gift *Gift) (*Gift, httperrors.HttpErr) {
	if err1 := Gift.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	Gift.Code = code
	Gift.Base.Updated_At = time.Now()
	Gift.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("Gift")
	result1, err := collection.InsertOne(ctx, &Gift)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create Gift Failed, %d", err))
	}
	Gift.ID = result1.InsertedID.(primitive.ObjectID)
	return Gift, nil
}
func (r giftrepository) GetOne(code string) (Gift *Gift, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("Gift")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&Gift)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return Gift, nil
}
func (r giftrepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.D{{"name", -1}})
	collection := db.Mongodb.Collection("Gift")
	results := []*Gift{}
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

func (r giftrepository) Update(code string, gift *Gift) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("Gift")
	var ac Gift
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if gift.Name == "" {
		gift.Name = ac.Name
	}
	if gift.Code == "" {
		gift.Code = ac.Code
	}
	if gift.Title == "" {
		gift.Title = ac.Title
	}
	if gift.Description == "" {
		gift.Description = ac.Description
	}
	update := bson.M{"$set": gift}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}

func (r giftrepository) GetOneByName(name string) (ac *Gift, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("Gift")
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

func (r giftrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("Gift")

	filter := bson.M{"code": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r giftrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("Gift")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "GiftCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r giftrepository) getuno(code string) (result *Gift, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("Gift")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}

func (r giftrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("Gift")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
