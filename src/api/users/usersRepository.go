package users

import (
	"context"
	"fmt"
	"strconv"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/db"
	"github.com/myrachanto/sports/src/pasetos"
	"github.com/myrachanto/sports/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Userrepository repository
var (
	Userrepository UserrepoInterface = &userrepository{}
	ctx                              = context.TODO()
	Userrepo                         = userrepository{}
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

// func LoadKey() (key Key, err error) {
// 	viper.AddConfigPath("../../../../")
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("env")

// 	viper.AutomaticEnv()

// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return
// 	}
// 	err = viper.Unmarshal(&key)
// 	return
// }

type UserrepoInterface interface {
	Create(user *User) (*User, httperrors.HttpErr)
	Login(auser *LoginUser) (*Auth, httperrors.HttpErr)
	RenewAccessToken(renewAccesstoken string) (*Auth, httperrors.HttpErr)
	Logout(token string) (string, httperrors.HttpErr)
	GetOne(id string) (*User, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	Forgot(email string) (string, string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Update(code string, user *User) httperrors.HttpErr
	PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr)
	PasswordReset(email, newpassword string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	UpdateAdmin(code string, status bool) httperrors.HttpErr
	UpdateAuditor(code string, status bool) httperrors.HttpErr
}
type userrepository struct{}

func NewUserRepo() UserrepoInterface {
	return &userrepository{}
}

func (r *userrepository) Create(user *User) (*User, httperrors.HttpErr) {

	if err1 := user.Validate(); err1 != nil {
		return nil, err1
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return nil, err1
	}
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return nil, httperrors.NewNotFoundError("Your email format is wrong!")
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	count, errd := r.Count()
	if errd != nil {
		return nil, errd
	}
	if count < 3 {
		user.Admin = true
	}

	user.Usercode = code
	user.Base.Updated_At = time.Now()
	user.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection("user")
	ok = r.emailexist(user.Email)
	if ok {
		return nil, httperrors.NewBadRequestError("that email exist in the our system!")
	}
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return nil, err2
	}
	user.Password = hashpassword

	result1, err := collection.InsertOne(ctx, &user)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create user Failed, %d", err))
	}
	user.ID = result1.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *userrepository) Login(user *LoginUser) (*Auth, httperrors.HttpErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	collection := db.Mongodb.Collection("user")
	filter := bson.M{"email": user.Email}
	auser := &User{}
	err := collection.FindOne(ctx, filter).Decode(&auser)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("User with this email does exist @ - , %d", err))
	}
	ok := user.Compare(user.Password, auser.Password)
	if !ok {
		return nil, httperrors.NewNotFoundError("wrong email password combo!")
	}
	maker, errs := pasetos.NewPasetoMaker()
	if err != nil {
		return nil, errs
	}
	tokencode, errs := Sessionsrepo.GeneTokencode(auser.Usercode)
	if errs != nil {
		return nil, errs
	}
	renewtokencode, errs := Sessionsrepo.GeneSessioncode(auser.Usercode)
	if errs != nil {
		return nil, errs
	}
	data := &pasetos.Data{
		Code:     tokencode,
		Usercode: auser.Usercode,
		Admin:    auser.Admin,
		Auditor:  auser.Auditor,
		Username: auser.Username,
		Email:    auser.Email,
	}
	// fmt.Println("---------------------", data)
	tokenString, payload, errs := maker.CreateToken(data, time.Hour*5)
	if errs != nil {
		return nil, errs
	}
	data.Code = renewtokencode
	RefleshToken, refleshtok, errs := maker.CreateToken(data, time.Hour*24)
	if errs != nil {
		return nil, errs
	}
	sessiond, errs := Sessionsrepo.CreateSession(&Session{
		Code: renewtokencode,
		// TokenId:      tokencode,
		Username:     auser.Username,
		Usercode:     auser.Usercode,
		RefleshToken: RefleshToken,
		UserAgent:    user.UserAgent,
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refleshtok.ExpiredAt,
	})
	if errs != nil {
		return nil, errs
	}
	auths := &Auth{Usercode: auser.Usercode, Picture: auser.Picture, UserName: auser.Username, Admin: auser.Admin, Token: tokenString, RefleshToken: RefleshToken, SessionCode: sessiond.Code, TokenExpires: payload.ExpiredAt, RefleshTokenExpires: sessiond.ExpiresAt}
	return auths, nil
}

func (r *userrepository) RenewAccessToken(renewAccesstoken string) (*Auth, httperrors.HttpErr) {
	maker, err := pasetos.NewPasetoMaker()
	if err != nil {
		return nil, err
	}
	refleshpayload, err := maker.VerifyToken(renewAccesstoken)
	if err != nil {
		return nil, err
	}
	sessions, err := Sessionsrepo.GetOne(refleshpayload.Code)
	if err != nil {
		return nil, err
	}
	if sessions.IsBlocked {
		if err != nil {
			return nil, httperrors.NewAnuthorizedError("your Session is blocked")
		}
	}
	if sessions.Username != refleshpayload.Username {
		if err != nil {
			return nil, httperrors.NewAnuthorizedError("your Session is blocked -u")
		}
	}

	tokencode, errs := Sessionsrepo.GeneTokencode(sessions.Usercode)
	if errs != nil {
		return nil, errs
	}
	tokenString, payload, errs := maker.CreateToken(&pasetos.Data{
		Username: refleshpayload.Username,
		Code:     tokencode,
		Usercode: sessions.Usercode,
		Email:    refleshpayload.Email,
	}, time.Hour*1)
	if errs != nil {
		return nil, errs
	}
	auths := &Auth{Usercode: sessions.Usercode, UserName: sessions.Username, Token: tokenString, TokenExpires: payload.ExpiredAt, RefleshTokenExpires: sessions.ExpiresAt}
	return auths, nil
}
func (r *userrepository) Logout(token string) (string, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(token)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("auth")
	filter1 := bson.M{"token": token}
	_, err3 := collection.DeleteOne(ctx, filter1)
	if err3 != nil {
		return "", httperrors.NewBadRequestError("something went wrong login out!")
	}
	return "something went wrong login out!", nil
}
func (r *userrepository) GetOne(code string) (user *User, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("user")
	filter := bson.M{"usercode": code}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return user, nil
}
func (r *userrepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("user")
	results := []*User{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.D{{"fullname", -1}})
	if search.Search != "" {
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"fullname", primitive.Regex{Pattern: search.Search, Options: "i"}}},
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
func (r *userrepository) UpdateAdmin(code string, status bool) httperrors.HttpErr {
	var n User

	collection := db.Mongodb.Collection("user")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"usercode", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"admin", status}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	errd := r.UpdateAuditor(code, status)
	if errd != nil {
		return errd
	}
	return nil
}
func (r *userrepository) UpdateAuditor(code string, status bool) httperrors.HttpErr {
	var n User

	collection := db.Mongodb.Collection("user")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"usercode", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	// fmt.Println("-------------------------auditor", status)
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"auditor", status}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}

func (r *userrepository) Update(code string, user *User) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}
	uuser := &User{}

	// ok := user.ValidateEmail(user.Email)
	// if !ok {
	// 	return httperrors.NewNotFoundError("Your email format is wrong!")
	// }
	// fmt.Println(code)
	user.Base.Updated_At = time.Now()
	collection := db.Mongodb.Collection("user")
	filter := bson.M{"usercode": code}
	err := collection.FindOne(ctx, filter).Decode(&uuser)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	fmt.Println("-------------------step 1")

	if user.Fullname == "" {
		user.Fullname = uuser.Fullname
	}
	if user.Username == "" {
		user.Username = uuser.Username
	}
	if user.Phone == "" {
		user.Phone = uuser.Phone
	}
	if user.Address == "" {
		user.Address = uuser.Address
	}
	if user.Picture == "" {
		user.Picture = uuser.Picture
	}
	if user.Email == "" {
		user.Email = uuser.Email
	}
	if user.Password == "" {
		user.Password = uuser.Password
	}
	if user.Usercode == "" {
		user.Usercode = uuser.Usercode
	}
	user.Admin = uuser.Admin
	user.Base.Created_At = uuser.Base.Created_At
	user.Base.Updated_At = time.Now()
	update := bson.M{"$set": user}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}

func (r *userrepository) PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(oldpassword)
	if stringresults.Noerror() {
		return "", "", stringresults
	}
	stringresults2 := httperrors.ValidStringNotEmpty(email)
	if stringresults2.Noerror() {
		return "", "", stringresults2
	}
	stringresults3 := httperrors.ValidStringNotEmpty(newpassword)
	if stringresults3.Noerror() {
		return "", "", stringresults3
	}
	upay := &User{}

	fmt.Println(oldpassword, email, newpassword)
	collection := db.Mongodb.Collection("user")
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return "", "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this email, %d", err))
	}
	ok := upay.Compare(oldpassword, upay.Password)
	if !ok {
		return "", "", httperrors.NewNotFoundError("wrong password combo!")
	}
	newhashpassword, err2 := upay.HashPassword(newpassword)
	if err2 != nil {
		return "", "", err2
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"password", newhashpassword}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", "", httperrors.NewNotFoundError("Error updating!")
	}
	return email, newpassword, nil
}
func (r *userrepository) PasswordReset(email, password string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(password)
	if stringresults.Noerror() {
		return "", stringresults
	}
	stringresults2 := httperrors.ValidStringNotEmpty(email)
	if stringresults2.Noerror() {
		return "", stringresults2
	}
	upay := &User{}

	collection := db.Mongodb.Collection("user")
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"email", primitive.Regex{Pattern: email, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this email, %d", err))
	}
	newhashpassword, err2 := upay.HashPassword(password)
	if err2 != nil {
		return "", err2
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"password", newhashpassword}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return email, nil
}
func (r *userrepository) Forgot(email string) (string, string, httperrors.HttpErr) {
	stringresults2 := httperrors.ValidStringNotEmpty(email)
	if stringresults2.Noerror() {
		return "", "", stringresults2
	}
	upay := &User{}

	collection := db.Mongodb.Collection("user")
	// filter := bson.M{"email": email}

	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"email", primitive.Regex{Pattern: email, Options: "i"}}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return "", "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this email, %d", err))
	}
	pass := support.Generatepassword()
	hashpassword, err2 := upay.HashPassword(pass)
	if err2 != nil {
		return "", "", err2
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"password", hashpassword}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", "", httperrors.NewNotFoundError("Error updating!")
	}
	return email, pass, nil
}

func (r userrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("user")

	// idPrimitive, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return "", httperrors.NewNotFoundError("primitive issue")
	// }
	filter := bson.M{"usercode": id}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r userrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("user")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "UserCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r userrepository) getuno(code string) (result *User, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("user")
	filter := bson.M{"usercode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r userrepository) emailexist(email string) bool {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return stringresults.Noerror()
	}
	collection := db.Mongodb.Collection("user")
	result := &User{}
	filter := bson.M{"email": email}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	return err1 == nil
}
func (r userrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("user")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
func (userRepo userrepository) Cleaner(code string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("user")

	filter := bson.M{"usercode": code}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil
}
