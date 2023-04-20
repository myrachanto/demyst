package users

import (
	"fmt"
	"strconv"
	"time"

	"github.com/myrachanto/demyst/src/db"
	"github.com/myrachanto/demyst/src/support"
	httperrors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Sessionsrepo sessionsrepo

type sessionsrepo struct{}

// /Dealing with sessions
func (r *sessionsrepo) CreateSession(session *Session) (*Session, httperrors.HttpErr) {

	code, err1 := r.GeneSessioncode(session.Usercode)
	if err1 != nil {
		return nil, err1
	}
	session.Base.Updated_At = time.Now()
	session.Base.Created_At = time.Now()
	session.Code = code
	collection := db.Mongodb.Collection("session")
	result, err := collection.InsertOne(ctx, session)
	if err != nil {
		// fmt.Println("err -------", err)
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create product Failed, %d", err))
	}
	session.ID = result.InsertedID.(primitive.ObjectID)
	return session, nil
}
func (r sessionsrepo) GeneSessioncode(usercode string) (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("session")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "session-" + usercode + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r sessionsrepo) GeneTokencode(usercode string) (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp
	cod := "token-" + usercode + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r *sessionsrepo) GetOne(code string) (session *Session, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("session")
	filter := bson.M{"code": code}
	fmt.Println("-----------", code)
	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return session, nil
}
