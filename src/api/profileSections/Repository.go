package profilesections

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
)

// Profilerepository repository
var (
	ProfileSectionrepository ProfileSectionrepoInterface = &profileSectionrepository{}
	ctx                                                  = context.TODO()
	database                                             = "Profilesections"
	ProfileSectionrepo                                   = profileSectionrepository{}
)

type ProfileSectionrepoInterface interface {
	Create(section *ProfileSection) (*ProfileSection, httperrors.HttpErr)
	GetOne(code string) (section *ProfileSection, errors httperrors.HttpErr)
	Update(code string, section *ProfileSection) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
}

type profileSectionrepository struct{}

func NewProfileSectionRepo() ProfileSectionrepoInterface {
	return &profileSectionrepository{}
}

func (r *profileSectionrepository) Create(section *ProfileSection) (*ProfileSection, httperrors.HttpErr) {
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	section.Code = code
	section.Base.Updated_At = time.Now()
	section.Base.Created_At = time.Now()
	r.EnsureRightData(section)
	collection := db.Mongodb.Collection(database)
	result1, errd := collection.InsertOne(ctx, &section)
	if errd != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create Profile Failed, %d", errd))
	}
	section.ID = result1.InsertedID.(primitive.ObjectID)
	return section, nil
}
func (r *profileSectionrepository) GetOne(code string) (section *ProfileSection, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection(database)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&section)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this code, %d", err))
	}
	return section, nil
}

func (r *profileSectionrepository) Update(code string, section *ProfileSection) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}

	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection(database)
	var ac ProfileSection
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if section.Name == "" {
		section.Name = ac.Name
	}
	if section.Code == "" {
		section.Code = ac.Code
	}
	if section.Content == "" {
		section.Content = ac.Content
	}
	if section.Image == "" {
		section.Image = ac.Image
	}
	if section.Code == "" {
		section.Code = ac.Code
	}
	if section.Highlight == "" {
		section.Highlight = ac.Highlight
	}
	section.Base.Updated_At = time.Now()
	r.EnsureRightData(section)
	update := bson.M{"$set": section}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}
func (r profileSectionrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection(database)
	// fmt.Println("stepper 1--------------------")
	res, errs := r.getuno(id)
	if errs != nil {
		return "", errs
	}
	// fmt.Println("stepper 2--------------------", res)
	go support.Clean.Cleaner(res.Image)
	filter := bson.M{"code": id}
	// fmt.Println("stepper 3--------------------", res)
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {

		// fmt.Println("stepper 4 err--------------------", res)
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}

func (r profileSectionrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection(database)
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "ProfileCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r profileSectionrepository) getuno(code string) (result *ProfileSection, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection(database)
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}

	return result, nil
}
func (r profileSectionrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection(database)
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
func (r profileSectionrepository) EnsureRightData(res *ProfileSection) {
	if res.Image != "" {
		res.Content = ""
	} else {
		res.Image = ""
	}
}
