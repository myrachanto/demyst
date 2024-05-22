package profile

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/api/category"
	profilesections "github.com/myrachanto/estate/src/api/profileSections"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Profilerepository repository
var (
	Profilerepository ProfileRepoInterface = &profilerepository{}
	ctx                                    = context.TODO()
	Profilerepo                            = profilerepository{}
)

type ProfileRepoInterface interface {
	Create(Profile *Profile) (*Profile, httperrors.HttpErr)
	GetOne(id string) (*Profile, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, Profile *Profile) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *Profile, errors httperrors.HttpErr)
	UpdateExclusive(code string, status bool) httperrors.HttpErr
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateTrending(code string, status bool) httperrors.HttpErr
	GetOneByUrl(code string) (Profile *ByProfile, errors httperrors.HttpErr)
	GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr)
	CreateComment(comment *Comment) httperrors.HttpErr
	DeleteComment(code, Profilecode string) httperrors.HttpErr
}
type profilerepository struct{}

func NewProfileRepo() ProfileRepoInterface {
	return &profilerepository{}
}

func (r *profilerepository) Create(Profile *Profile) (*Profile, httperrors.HttpErr) {
	if err1 := Profile.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	cat, err := category.Categoryrepo.GetOneByName(Profile.Sport)
	if err != nil {
		return nil, err
	}
	Profile.Sport = cat.Code
	Profile.Code = code
	Profile.Base.Updated_At = time.Now()
	Profile.Base.Created_At = time.Now()
	ress := []string{}
	for _, g := range Profile.Sections {
		res, _ := profilesections.ProfileSectionrepo.Create(g)
		ress = append(ress, res.Code)
	}
	Profile.SectionsCodes = ress
	collection := db.Mongodb.Collection("Profile")
	result1, errd := collection.InsertOne(ctx, &Profile)
	if errd != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create Profile Failed, %d", errd))
	}
	Profile.ID = result1.InsertedID.(primitive.ObjectID)
	return r.ModifyOneProfileWithCategory(Profile), nil
}
func (r *profilerepository) CreateComment(comment *Comment) httperrors.HttpErr {
	Profile, err := r.getuno(comment.Profilecode)
	comments := Profile.Comments
	if err != nil {
		return err
	}
	comment.Code = uuid.New().String()
	comment.Base.Updated_At = time.Now()
	comment.Base.Created_At = time.Now()
	comments = append(comments, *comment)
	collection := db.Mongodb.Collection("Profile")
	fmt.Println("------------------")
	filter := bson.M{"code": comment.Profilecode}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"comments", comments}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *profilerepository) DeleteComment(code, Profilecode string) httperrors.HttpErr {
	Profile, err := r.getuno(Profilecode)
	if err != nil {
		return err
	}
	comments := []*Comment{}
	for _, g := range Profile.Comments {
		if g.Code != code {
			comments = append(comments, &g)
		}
	}
	collection := db.Mongodb.Collection("Profile")
	fmt.Println("------------------")
	filter := bson.M{"code": Profilecode}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"comments", comments}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *profilerepository) GetOne(code string) (profile *Profile, errors httperrors.HttpErr) {
	var n Profile
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("Profile")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this code, %d", err))
	}
	return r.GetProfileections(&n), nil
}
func (r *profilerepository) GetOneByUrl(code string) (profile *ByProfile, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("Profile")
	filter := bson.M{"url": code}
	var report Profile
	err := collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this url, %d", err))
	}
	trending, e := r.GetAllTrending(20)
	if e != nil {
		return nil, e
	}
	rep := r.GetProfileections(&report)
	return &ByProfile{
		Profile:  rep,
		Trending: trending,
	}, nil
}

func (r *profilerepository) GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr) {
	// fmt.Println("=++++++++++++++++++++++++step1sxascs ")
	collection := db.Mongodb.Collection("Profile")
	results := []*Profile{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.M{"base.updated_at": -1})
	res, errs := category.Categoryrepo.GetOneByName(code)
	// fmt.Println("=++++++++++++++++++++++++step3sxascs ", res, errs)
	if errs != nil {
		return nil, errs
	}
	// fmt.Println("=++++++++++++++++++++++++step3sxascs ", res.Code)

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"sport", res.Code}},
		}},
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// fmt.Println("=++++++++++++++++++++++++step2", len(results))
	final := []*Profile{}
	for _, m := range results {
		res := r.GetProfileections(m)
		final = append(final, res)
	}
	count, err := collection.CountDocuments(ctx, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	return &Results{
		Data:  final,
		Total: int(count),
	}, nil

}
func (r *profilerepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	count, errd := r.Count()
	if errd != nil {
		return nil, errd
	}
	collection := db.Mongodb.Collection("Profile")
	results := []*Profile{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.M{"base.updated_at": -1})
	if search.Search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search.Search, Options: "i"}}},
				bson.D{{"title", primitive.Regex{Pattern: search.Search, Options: "i"}}},
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

		count, err := collection.CountDocuments(ctx, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		final := []*Profile{}
		for _, v := range results {
			res := r.GetProfileections(v)
			final = append(final, res)
		}
		return &Results{
			Data:  final,
			Total: int(count),
		}, nil
	} else {
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		final := []*Profile{}
		for _, v := range results {
			res := r.GetProfileections(v)
			final = append(final, res)
		}
		return &Results{
			Data:  final,
			Total: int(count),
		}, nil
	}

}
func (r *profilerepository) GetAll1() ([]*Profile, httperrors.HttpErr) {
	results := []*Profile{}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"base.updated_at": -1})
	collection := db.Mongodb.Collection("Profile")
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*Profile{}
	for _, v := range results {
		res := r.GetProfileections(v)
		final = append(final, res)
	}
	return final, nil

}
func (r *profilerepository) GetAll2(limit int64) ([]*Profile, httperrors.HttpErr) {
	results := []*Profile{}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"base.updated_at": -1})
	findOptions.SetLimit(limit)
	collection := db.Mongodb.Collection("Profile")
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*Profile{}
	for _, v := range results {
		res := r.GetProfileections(v)
		final = append(final, res)
	}
	return final, nil

}

func (r *profilerepository) GetAllTrending(num ...int64) ([]*Profile, httperrors.HttpErr) {
	results := []*Profile{}
	collection := db.Mongodb.Collection("Profile")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"trending", true}},
		}},
	}
	var limit int64
	if len(num) == 0 {
		limit = 1000
	} else {
		limit = num[0]
	}
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.M{"base.updated_at": -1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*Profile{}
	for _, v := range results {
		res := r.GetProfileections(v)
		final = append(final, res)
	}
	return final, nil

}

func (r *profilerepository) GetAllExclusive(num ...int64) ([]*Profile, httperrors.HttpErr) {
	results := []*Profile{}
	collection := db.Mongodb.Collection("Profile")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"exclusive", true}},
		}},
	}
	var limit int64
	if len(num) == 0 {
		limit = 1000
	} else {
		limit = num[0]
	}
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.M{"base.updated_at": -1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*Profile{}
	for _, v := range results {
		res := r.GetProfileections(v)
		final = append(final, res)
	}
	return final, nil

}
func (r *profilerepository) GetAllCategory() ([]*SportCount, httperrors.HttpErr) {
	results := []*SportCount{}
	collection := db.Mongodb.Collection("Profile")
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$sport"},
			{"count", bson.D{
				{"$sum", 1},
			}},
		}},
	}
	pipeline := mongo.Pipeline{groupStage}
	// cursor, err := collection.Find(ctx, filter, opts)
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, n := range results {
		cat, _ := category.Categoryrepo.GetOne(n.Sport)
		n.Sport = cat.Name
	}
	return results, nil

}

func (r *profilerepository) GetAllPostByWeek() ([]*Weekly, httperrors.HttpErr) {
	results := []*Weekly{}
	collection := db.Mongodb.Collection("Profile")
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"week": bson.M{
					"$week": "$base.created_at",
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$week",
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	// Aggregation pipeline
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	for cursor.Next(context.TODO()) {

		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}

		week := result["_id"].(int32)
		count := result["count"].(int32)

		// fmt.Printf("Week %d: %d documents\n", week, count)
		results = append(results, &Weekly{week, count})
	}

	if err := cursor.Err(); err != nil {
		return nil, httperrors.NewNotFoundError("Cursor error!")
	}

	return results, nil

}
func (r *profilerepository) GetAllFeatured(num ...int64) ([]*Profile, httperrors.HttpErr) {
	results := []*Profile{}
	collection := db.Mongodb.Collection("Profile")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"featured", true}},
		}},
	}
	var limit int64
	if len(num) == 0 {
		limit = 1000
	} else {
		limit = num[0]
	}
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.M{"base.updated_at": -1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*Profile{}
	for _, v := range results {
		res := r.GetProfileections(v)
		final = append(final, res)
	}
	return final, nil

}

func (r *profilerepository) Update(code string, profile *Profile) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	if profile.Sport != "" {
		cat, err := category.Categoryrepo.GetOneByName(profile.Sport)
		if err != nil {
			return "", err
		}
		profile.Sport = cat.Code
	}
	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("Profile")
	var ac Profile
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if profile.Name == "" {
		profile.Name = ac.Name
	}
	if profile.Code == "" {
		profile.Code = ac.Code
	}
	if profile.Title == "" {
		profile.Title = ac.Title
	}
	if profile.Content == "" {
		profile.Content = ac.Content
	}
	if profile.Picture == "" {
		profile.Picture = ac.Picture
	}
	if profile.Meta == "" {
		profile.Meta = ac.Meta
	}
	if profile.Caption == "" {
		profile.Caption = ac.Caption
	}
	if profile.Url == "" {
		profile.Url = ac.Url
	}
	if profile.Sport == "" {
		profile.Sport = ac.Sport
	}
	if !profile.Featured {
		profile.Featured = ac.Featured
	}
	if !profile.Exclusive {
		profile.Exclusive = ac.Exclusive
	}
	if !profile.Trending {
		profile.Trending = ac.Trending
	}
	profile.Base.Created_At = ac.Base.Created_At
	ress := []string{}
	for _, sec := range profile.Sections {
		if sec.Code == "" {
			// Create a new secion
			new, errs := profilesections.ProfileSectionrepo.Create(sec)
			if errs != nil {
				fmt.Println(errs.Message())
			}
			ress = append(ress, new.Code)
		} else {
			// update a new sections
			profilesections.ProfileSectionrepo.Update(sec.Code, sec)
			ress = append(ress, sec.Code)
		}
	}
	profile.SectionsCodes = ress
	fmt.Println("----------------step 1a", profile.SectionsCodes)
	//delete items
	for _, gv := range ac.SectionsCodes {
		fmt.Println("----------------step 1", gv)
		coder, ok := r.EvaluateIfSecionExist(gv, profile.Sections)
		if !ok {
			fmt.Println("---------------- step 2", gv, coder)
			profilesections.ProfileSectionrepo.Delete(coder)
		}
	}
	profile.Base.Updated_At = time.Now()
	update := bson.M{"$set": profile}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}
func (r *profilerepository) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	var n Profile

	collection := db.Mongodb.Collection("Profile")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
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
			{"$set", bson.D{{"featured", status}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}

func (r *profilerepository) UpdateTrending(code string, status bool) httperrors.HttpErr {
	var n Profile

	collection := db.Mongodb.Collection("Profile")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	// fmt.Println("++++++++++++++++++", code, status)
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"trending", status}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *profilerepository) UpdateExclusive(code string, status bool) httperrors.HttpErr {
	var n Profile

	collection := db.Mongodb.Collection("Profile")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
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
			{"$set", bson.D{{"exclusive", status}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}

func (r *profilerepository) GetOneByName(name string) (acc *Profile, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("Profile")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}},
		}},
	}
	var ac Profile
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	Profile := r.GetProfileections(&ac)
	return Profile, nil
}
func (r profilerepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)

	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("Profile")

	filter := bson.M{"code": id}

	res, err := r.getuno(id)
	if err != nil {
		return "", err
	}
	go support.Clean.Cleaner(res.Picture)
	for _, m := range res.Sections {
		_, _ = profilesections.ProfileSectionrepo.Delete(m.Code)
	}
	ok, errs := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", errs))
	}
	return "deleted successfully", nil

}
func (r profilerepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("Profile")
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
func (r profilerepository) getuno(code string) (result *Profile, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	var res Profile
	collection := db.Mongodb.Collection("Profile")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&res)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}

	return r.GetProfileections(&res), nil
}
func (r profilerepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("Profile")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}

func (r profilerepository) ModifyOneProfileWithCategory(Profile *Profile) *Profile {
	cat, _ := category.Categoryrepo.GetOne(Profile.Sport)
	Profile.Sport = cat.Name
	return Profile
}
func (r profilerepository) ModifyProfileWithCategory(Profile []*Profile) []*Profile {
	for _, n := range Profile {
		cat, _ := category.Categoryrepo.GetOne(n.Sport)
		n.Sport = cat.Name
	}
	return Profile
}

//	func (r Profilerepository) ConvertoProfileB(Profile *Profile) *ProfileB {
//		var n ProfileB
//		n.Name = Profile.Name
//		n.Title = Profile.Title
//		n.Caption = Profile.Caption
//		n.Meta = Profile.Meta
//		n.Url = Profile.Url
//		n.Sport = Profile.Sport
//		n.Featured = Profile.Featured
//		n.Exclusive = Profile.Exclusive
//		n.Trending = Profile.Trending
//		n.Content = Profile.Content
//		n.Picture = Profile.Picture
//		n.Code = Profile.Code
//		n.Base = Profile.Base
//		n.Comments = Profile.Comments
//		//sections
//		results := []Coder{}
//		for _, g := range Profile.Sections {
//			results = append(results, Coder{Name: g.Code})
//		}
//		n.Sections = results
//		return &n
//	}
func (r profilerepository) GetProfileections(Profile *Profile) *Profile {
	for _, m := range Profile.SectionsCodes {
		res, _ := profilesections.ProfileSectionrepo.GetOne(m)
		Profile.Sections = append(Profile.Sections, res)
	}
	r.ModifyOneProfileWithCategory(Profile)
	return Profile
}
func (r profilerepository) EvaluateIfSecionExist(code string, oldinfo []*profilesections.ProfileSection) (string, bool) {
	var results bool = false
	for _, g := range oldinfo {
		if g.Code == code {
			results = true
			break
		}
	}
	return code, results
}
func (r profilerepository) GetIfSecionExist(code string, newInfo []*profilesections.ProfileSection) *profilesections.ProfileSection {
	var results profilesections.ProfileSection
	for _, g := range newInfo {
		if g.Code == code {
			results = *g
			break
		}
	}
	return &results
}
func (r profilerepository) DeleteAsection(code string, secs []Coder) []Coder {
	l := len(secs) - 1
	Profilelice := make([]Coder, l, l)
	for _, g := range secs {
		if g.Name != code {
			Profilelice = append(Profilelice, g)
		}
	}
	return Profilelice
}

func EvaluateIfSecionExist(code int, newInfo []int) (int, bool) {
	var results bool = false
	for _, g := range newInfo {
		if g == code {
			results = true
			break
		}
	}
	return code, results
}
