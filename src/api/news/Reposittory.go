package news

import (
	"context"
	"fmt"
	"strconv"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/api/category"
	newssections "github.com/myrachanto/sports/src/api/newsSections"
	"github.com/myrachanto/sports/src/db"
	"github.com/myrachanto/sports/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// newsrepository repository
var (
	Newsrepository NewsrepoInterface = &newsrepository{}
	ctx                              = context.TODO()
	Newsrepo                         = newsrepository{}
)

type NewsrepoInterface interface {
	Create(news *News) (*News, httperrors.HttpErr)
	GetOne(id string) (*News, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, news *News) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *News, errors httperrors.HttpErr)
	UpdateExclusive(code string, status bool) httperrors.HttpErr
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateTrending(code string, status bool) httperrors.HttpErr
	GetOneByUrl(code string) (news *ByNews, errors httperrors.HttpErr)
	GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr)
}
type newsrepository struct{}

func NewnewsRepo() NewsrepoInterface {
	return &newsrepository{}
}

func (r *newsrepository) Create(news *News) (*News, httperrors.HttpErr) {
	if err1 := news.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	cat, err := category.Categoryrepo.GetOneByName(news.Sport)
	if err != nil {
		return nil, err
	}
	news.Sport = cat.Code
	news.Code = code
	news.Base.Updated_At = time.Now()
	news.Base.Created_At = time.Now()
	res := r.ConvertoNewsB(news)
	collection := db.Mongodb.Collection("news")
	result1, errd := collection.InsertOne(ctx, &res)
	if errd != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create news Failed, %d", errd))
	}
	res.ID = result1.InsertedID.(primitive.ObjectID)
	news = r.ConvertoNews(res)
	return news, nil
}
func (r *newsrepository) GetOne(code string) (news *News, errors httperrors.HttpErr) {
	var n NewsB
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("news")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this code, %d", err))
	}
	return r.ModifyOneNewsWithCategory(r.ConvertoNews(&n)), nil
}
func (r *newsrepository) GetOneByUrl(code string) (news *ByNews, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("news")
	filter := bson.M{"url": code}
	var report NewsB
	err := collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this url, %d", err))
	}
	trending, e := r.GetAllTrending()
	if e != nil {
		return nil, e
	}
	n := r.ConvertoNews(&report)
	rep := r.ModifyOneNewsWithCategory(n)
	return &ByNews{
		News:     rep,
		Trending: trending,
	}, nil
}

func (r *newsrepository) GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr) {
	// fmt.Println("=++++++++++++++++++++++++step1sxascs ")
	collection := db.Mongodb.Collection("news")
	results := []*NewsB{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.D{{"name", -1}})
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
	final := []*News{}
	for _, m := range results {
		res := r.ConvertoNews(m)
		final = append(final, res)
	}
	count, err := collection.CountDocuments(ctx, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	fmt.Println("=++++++++++++++++++++++++step2", len(final))
	finals := r.ModifyNewsWithCategory(final)
	return &Results{
		Data:  finals,
		Total: int(count),
	}, nil

}
func (r *newsrepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	count, errd := r.Count()
	if errd != nil {
		return nil, errd
	}
	collection := db.Mongodb.Collection("news")
	results := []*NewsB{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	findOptions.SetSort(bson.D{{"name", -1}})
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
		final := []*News{}
		for _, v := range results {
			res := r.ConvertoNews(v)
			final = append(final, res)
		}
		finals := r.ModifyNewsWithCategory(final)
		return &Results{
			Data:  finals,
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
		final := []*News{}
		for _, v := range results {
			res := r.ConvertoNews(v)
			final = append(final, res)
		}
		finals := r.ModifyNewsWithCategory(final)
		return &Results{
			Data:  finals,
			Total: int(count),
		}, nil
	}

}
func (r *newsrepository) GetAll1() ([]*News, httperrors.HttpErr) {
	results := []*NewsB{}
	collection := db.Mongodb.Collection("news")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*News{}
	for _, v := range results {
		res := r.ConvertoNews(v)
		final = append(final, res)
	}
	finals := r.ModifyNewsWithCategory(final)
	return finals, nil

}

func (r *newsrepository) GetAllTrending(num ...int64) ([]*News, httperrors.HttpErr) {
	results := []*NewsB{}
	collection := db.Mongodb.Collection("news")
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
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*News{}
	for _, v := range results {
		res := r.ConvertoNews(v)
		final = append(final, res)
	}
	finals := r.ModifyNewsWithCategory(final)
	return finals, nil

}

func (r *newsrepository) GetAllExclusive(num ...int64) ([]*News, httperrors.HttpErr) {
	results := []*NewsB{}
	collection := db.Mongodb.Collection("news")
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
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*News{}
	for _, v := range results {
		res := r.ConvertoNews(v)
		final = append(final, res)
	}
	finals := r.ModifyNewsWithCategory(final)
	return finals, nil

}
func (r *newsrepository) GetAllCategory() ([]*SportCount, httperrors.HttpErr) {
	results := []*SportCount{}
	collection := db.Mongodb.Collection("news")
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

func (r *newsrepository) GetAllPostByWeek() ([]*Weekly, httperrors.HttpErr) {
	results := []*Weekly{}
	collection := db.Mongodb.Collection("news")
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
func (r *newsrepository) GetAllFeatured(num ...int64) ([]*News, httperrors.HttpErr) {
	results := []*NewsB{}
	collection := db.Mongodb.Collection("news")
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
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	final := []*News{}
	for _, v := range results {
		res := r.ConvertoNews(v)
		final = append(final, res)
	}
	finals := r.ModifyNewsWithCategory(final)
	return finals, nil

}

func (r *newsrepository) Update(code string, news *News) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	if news.Sport != "" {
		cat, err := category.Categoryrepo.GetOneByName(news.Sport)
		if err != nil {
			return "", err
		}
		news.Sport = cat.Code
	}
	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("news")
	var ac NewsB
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if news.Name == "" {
		news.Name = ac.Name
	}
	if news.Code == "" {
		news.Code = ac.Code
	}
	if news.Title == "" {
		news.Title = ac.Title
	}
	if news.Content == "" {
		news.Content = ac.Content
	}
	if news.Picture == "" {
		news.Picture = ac.Picture
	}
	if news.Meta == "" {
		news.Meta = ac.Meta
	}
	if news.Caption == "" {
		news.Caption = ac.Caption
	}
	if news.Url == "" {
		news.Url = ac.Url
	}
	if news.Sport == "" {
		news.Sport = ac.Sport
	}
	if !news.Featured {
		news.Featured = ac.Featured
	}
	if !news.Exclusive {
		news.Exclusive = ac.Exclusive
	}
	if !news.Trending {
		news.Trending = ac.Trending
	}
	news.Base.Created_At = ac.Base.Created_At
	for _, sec := range news.Sections {
		_, err := newssections.NewsSectionrepo.GetOne(sec.Code)
		if err != nil {
			// Create a new secion
			s, _ := newssections.NewsSectionrepo.Create(&sec)
			news.Sections = append(news.Sections, *s)
		} else {
			// update a new sections
			newssections.NewsSectionrepo.Update(sec.Code, &sec)
		}
	}
	// delete a section
	res := r.ConvertoNewsB(news)
	for _, gv := range ac.Sections {
		coder, ok := r.EvaluateIfSecionExist(gv.Name, news.Sections)
		if !ok {
			newssections.NewsSectionrepo.Delete(coder)
			res.Sections = r.DeleteAsection(coder, ac.Sections)
		}
	}
	res.Base.Updated_At = time.Now()
	update := bson.M{"$set": res}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}
func (r *newsrepository) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	var n News

	collection := db.Mongodb.Collection("news")

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

func (r *newsrepository) UpdateTrending(code string, status bool) httperrors.HttpErr {
	var n News

	collection := db.Mongodb.Collection("news")

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
func (r *newsrepository) UpdateExclusive(code string, status bool) httperrors.HttpErr {
	var n News

	collection := db.Mongodb.Collection("news")

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

func (r *newsrepository) GetOneByName(name string) (acc *News, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("news")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}},
		}},
	}
	var ac NewsB
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	news := r.ConvertoNews(&ac)
	return r.ModifyOneNewsWithCategory(news), nil
}
func (r newsrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)

	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("news")

	filter := bson.M{"code": id}

	res, err := r.getuno(id)
	if err != nil {
		return "", err
	}
	for _, m := range res.Sections {
		_, _ = newssections.NewsSectionrepo.Delete(m.Code)
	}
	ok, errs := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", errs))
	}
	return "deleted successfully", nil

}
func (r newsrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("news")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "newsCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r newsrepository) getuno(code string) (result *News, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	var res NewsB
	collection := db.Mongodb.Collection("news")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&res)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}

	return r.ModifyOneNewsWithCategory(r.ConvertoNews(&res)), nil
}
func (r newsrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("news")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}
func (r newsrepository) ModifyOneNewsWithCategory(news *News) *News {
	cat, _ := category.Categoryrepo.GetOne(news.Sport)
	news.Sport = cat.Name
	return news
}
func (r newsrepository) ModifyNewsWithCategory(news []*News) []*News {
	for _, n := range news {
		cat, _ := category.Categoryrepo.GetOne(n.Sport)
		n.Sport = cat.Name
	}
	return news
}

func (r newsrepository) ConvertoNewsB(news *News) *NewsB {
	var n NewsB
	n.Name = news.Name
	n.Title = news.Title
	n.Caption = news.Caption
	n.Meta = news.Meta
	n.Url = news.Url
	n.Sport = news.Sport
	n.Featured = news.Featured
	n.Exclusive = news.Exclusive
	n.Trending = news.Trending
	n.Content = news.Content
	n.Picture = news.Picture
	n.Code = news.Code
	n.Base = news.Base
	n.Comments = news.Comments
	//sections
	results := []Coder{}
	for _, g := range news.Sections {
		res, _ := newssections.NewsSectionrepo.Create(&g)
		results = append(results, Coder{Name: res.Code})
	}
	n.Sections = results
	return &n
}
func (r newsrepository) ConvertoNews(news *NewsB) *News {
	var n News
	n.Name = news.Name
	n.Title = news.Title
	n.Caption = news.Caption
	n.Meta = news.Meta
	n.Url = news.Url
	n.Sport = news.Sport
	n.Featured = news.Featured
	n.Exclusive = news.Exclusive
	n.Trending = news.Trending
	n.Content = news.Content
	n.Picture = news.Picture
	n.Code = news.Code
	n.Base = news.Base
	n.Comments = news.Comments
	//sections
	results := []newssections.NewsSection{}
	for _, m := range news.Sections {
		res, _ := newssections.NewsSectionrepo.GetOne(m.Name)
		results = append(results, *res)
	}
	n.Sections = results
	return &n
}
func (r newsrepository) EvaluateIfSecionExist(code string, newInfo []newssections.NewsSection) (string, bool) {
	var results bool = false
	for _, g := range newInfo {
		if g.Code == code {
			results = true
			break
		}
	}
	return code, results
}
func (r newsrepository) DeleteAsection(code string, secs []Coder) []Coder {
	l := len(secs) - 1
	newslice := make([]Coder, l, l)
	for _, g := range secs {
		if g.Name != code {
			newslice = append(newslice, g)
		}
	}
	return newslice
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
