package blog

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	httperrors "github.com/myrachanto/erroring"

	// blogsections "github.com/myrachanto/estate/src/api/BlogSections"
	"github.com/myrachanto/estate/src/api/category"
	"github.com/myrachanto/estate/src/api/product"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Blogrepository repository
var (
	Blogrepository BlogRepoInterface = &blogrepository{}
	ctx                              = context.TODO()
	Blogrepo                         = blogrepository{}
)

type BlogRepoInterface interface {
	Create(Blog *Blog) (*Blog, httperrors.HttpErr)
	GetOne(id string) (*Blog, httperrors.HttpErr)
	GetAll(support.Paginator) (*Results, httperrors.HttpErr)
	Update(code string, Blog *Blog) (string, httperrors.HttpErr)
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	GetOneByName(name string) (ac *Blog, errors httperrors.HttpErr)
	UpdateExclusive(code string, status bool) httperrors.HttpErr
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateTrending(code string, status bool) httperrors.HttpErr
	GetOneByUrl(code string) (Blog *ByBlog, errors httperrors.HttpErr)
	GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr)
	CreateComment(comment *Comment) httperrors.HttpErr
	DeleteComment(code, Blogcode string) httperrors.HttpErr
}
type blogrepository struct{}

func NewBlogRepo() BlogRepoInterface {
	return &blogrepository{}
}

func (r *blogrepository) Create(blog *Blog) (*Blog, httperrors.HttpErr) {
	fmt.Println("------------------step1", blog.Title)
	if err1 := blog.Validate(); err1 != nil {
		return nil, err1
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	// cat, err := category.Categoryrepo.GetOneByName(blog.Sport)
	// if err != nil {
	// 	return nil, err
	// }
	// blog.Sport = cat.Code
	blog.Code = code
	blog.Base.Updated_At = time.Now()
	blog.Base.Created_At = time.Now()
	// ress := []string{}
	// for _, g := range Blog.Sections {
	// 	res, _ := blogsections.BlogSectionrepository.Create(g)
	// 	ress = append(ress, res.Code)
	// }
	// Blog.SectionsCodes = ress
	collection := db.Mongodb.Collection("blog")
	result1, errd := collection.InsertOne(ctx, &blog)
	if errd != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create Blog Failed, %d", errd))
	}
	blog.ID = result1.InsertedID.(primitive.ObjectID)
	return blog, nil
}
func (r *blogrepository) CreateComment(comment *Comment) httperrors.HttpErr {
	Blog, err := r.getuno(comment.Code)
	comments := Blog.Comments
	if err != nil {
		return err
	}
	comment.Code = uuid.New().String()
	comment.Base.Updated_At = time.Now()
	comment.Base.Created_At = time.Now()
	comments = append(comments, *comment)
	collection := db.Mongodb.Collection("blog")
	fmt.Println("------------------")
	filter := bson.M{"code": comment.Code}
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
func (r *blogrepository) DeleteComment(code, Blogcode string) httperrors.HttpErr {
	Blog, err := r.getuno(Blogcode)
	if err != nil {
		return err
	}
	comments := []*Comment{}
	for _, g := range Blog.Comments {
		if g.Code != code {
			comments = append(comments, &g)
		}
	}
	collection := db.Mongodb.Collection("blog")
	fmt.Println("------------------")
	filter := bson.M{"code": Blogcode}
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

func (r *blogrepository) GetOne(code string) (blog *Blog, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("blog")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return blog, nil
}
func (r *blogrepository) GetOneByUrl(code string) (blog *ByBlog, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("blog")
	filter := bson.M{"url": code}
	var report Blog
	err := collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this url, %d", err))
	}
	trending, e := product.Productrepository.GetFeatured()
	if e != nil {
		return nil, e
	}
	// rep := r.GetBlogections(&report)
	return &ByBlog{
		Blog:     &report,
		Trending: trending,
	}, nil
}

func (r *blogrepository) GetByCategory(code string, search support.Paginator) (*Results, httperrors.HttpErr) {
	// fmt.Println("=++++++++++++++++++++++++step1sxascs ")
	collection := db.Mongodb.Collection("blog")
	results := []*Blog{}
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
	// final := []*Blog{}
	// for _, m := range results {
	// 	res := r.GetBlogections(m)
	// 	final = append(final, res)
	// }
	count, err := collection.CountDocuments(ctx, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	return &Results{
		Data:  results,
		Total: int(count),
	}, nil

}
func (r *blogrepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	count, errd := r.Count()
	if errd != nil {
		return nil, errd
	}
	collection := db.Mongodb.Collection("blog")
	results := []*Blog{}
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
		// final := []*Blog{}
		// for _, v := range results {
		// 	res := r.GetBlogections(v)
		// 	final = append(final, res)
		// }
		return &Results{
			Data:  results,
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
		// final := []*Blog{}
		// for _, v := range results {
		// 	res := r.GetBlogections(v)
		// 	final = append(final, res)
		// }
		return &Results{
			Data:  results,
			Total: int(count),
		}, nil
	}

}
func (r *blogrepository) GetAll1() ([]*Blog, httperrors.HttpErr) {
	results := []*Blog{}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"base.updated_at": -1})
	collection := db.Mongodb.Collection("blog")
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// final := []*Blog{}
	// for _, v := range results {
	// 	res := r.GetBlogections(v)
	// 	final = append(final, res)
	// }
	return results, nil

}
func (r *blogrepository) GetAll2(limit int64) ([]*Blog, httperrors.HttpErr) {
	results := []*Blog{}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"base.updated_at": -1})
	findOptions.SetLimit(limit)
	collection := db.Mongodb.Collection("blog")
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// final := []*Blog{}
	// for _, v := range results {
	// 	res := r.GetBlogections(v)
	// 	final = append(final, res)
	// }
	return results, nil

}

func (r *blogrepository) GetAllTrending(num ...int64) ([]*Blog, httperrors.HttpErr) {
	results := []*Blog{}
	collection := db.Mongodb.Collection("blog")
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
	// final := []*Blog{}
	// for _, v := range results {
	// 	res := r.GetBlogections(v)
	// 	final = append(final, res)
	// }
	return results, nil

}

func (r *blogrepository) GetAllExclusive(num ...int64) ([]*Blog, httperrors.HttpErr) {
	results := []*Blog{}
	collection := db.Mongodb.Collection("blog")
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
	// final := []*Blog{}
	// for _, v := range results {
	// 	res := r.GetBlogections(v)
	// 	final = append(final, res)
	// }
	return results, nil

}
func (r *blogrepository) GetAllCategory() ([]*SportCount, httperrors.HttpErr) {
	results := []*SportCount{}
	collection := db.Mongodb.Collection("blog")
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

func (r *blogrepository) GetAllPostByWeek() ([]*Weekly, httperrors.HttpErr) {
	results := []*Weekly{}
	collection := db.Mongodb.Collection("blog")
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
func (r *blogrepository) GetAllFeatured(num ...int64) ([]*Blog, httperrors.HttpErr) {
	results := []*Blog{}
	collection := db.Mongodb.Collection("blog")
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
	// final := []*Blog{}
	// for _, v := range results {
	// 	res := r.GetBlogections(v)
	// 	final = append(final, res)
	// }
	return results, nil

}

func (r *blogrepository) Update(code string, blog *Blog) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	// if blog.Sport != "" {
	// 	cat, err := category.Categoryrepo.GetOneByName(blog.Sport)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	blog.Sport = cat.Code
	// }
	filter := bson.M{"code": code}
	collection := db.Mongodb.Collection("blog")
	var ac Blog
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	if blog.Name == "" {
		blog.Name = ac.Name
	}
	if blog.Code == "" {
		blog.Code = ac.Code
	}
	if blog.Title == "" {
		blog.Title = ac.Title
	}
	if blog.Content == "" {
		blog.Content = ac.Content
	}
	if blog.Picture == "" {
		blog.Picture = ac.Picture
	}
	if blog.Meta == "" {
		blog.Meta = ac.Meta
	}
	if blog.Caption == "" {
		blog.Caption = ac.Caption
	}
	if blog.Url == "" {
		blog.Url = ac.Url
	}
	blog.Base.Created_At = ac.Base.Created_At
	blog.Base.Updated_At = time.Now()
	update := bson.M{"$set": blog}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "successifully Updated!", nil
}
func (r *blogrepository) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	var n Blog

	collection := db.Mongodb.Collection("blog")

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

func (r *blogrepository) UpdateTrending(code string, status bool) httperrors.HttpErr {
	var n Blog

	collection := db.Mongodb.Collection("blog")

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
func (r *blogrepository) UpdateExclusive(code string, status bool) httperrors.HttpErr {
	var n Blog

	collection := db.Mongodb.Collection("Blog")

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

func (r *blogrepository) GetOneByName(name string) (acc *Blog, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("blog")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}},
		}},
	}
	var ac Blog
	err := collection.FindOne(ctx, filter).Decode(&ac)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	// Blog :=
	return &ac, nil
}
func (r blogrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)

	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("blog")

	filter := bson.M{"code": id}

	res, err := r.getuno(id)
	if err != nil {
		return "", err
	}
	go support.Clean.Cleaner(res.Picture)
	// for _, m := range res.Sections {
	// 	_, _ = blogsections.BlogSectionrepo.Delete(m.Code)
	// }
	ok, errs := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", errs))
	}
	return "deleted successfully", nil

}
func (r blogrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection("blog")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "BlogCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r blogrepository) getuno(code string) (result *Blog, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	var res Blog
	collection := db.Mongodb.Collection("blog")
	filter := bson.M{"code": code}
	err1 := collection.FindOne(ctx, filter).Decode(&res)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}

	return &res, nil
}
func (r blogrepository) Count() (float64, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("blog")
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, httperrors.NewNotFoundError("no results found")
	}
	code := float64(count)
	return code, nil
}

// func (r blogrepository) ModifyOneBlogWithCategory(Blog *Blog) *Blog {
// 	cat, _ := category.Categoryrepo.GetOne(Blog.Sport)
// 	Blog.Sport = cat.Name
// 	return Blog
// }
// func (r blogrepository) ModifyBlogWithCategory(Blog []*Blog) []*Blog {
// 	for _, n := range Blog {
// 		cat, _ := category.Categoryrepo.GetOne(n.Sport)
// 		n.Sport = cat.Name
// 	}
// 	return Blog
// }

//	func (r Blogrepository) ConvertoBlogB(Blog *Blog) *BlogB {
//		var n BlogB
//		n.Name = Blog.Name
//		n.Title = Blog.Title
//		n.Caption = Blog.Caption
//		n.Meta = Blog.Meta
//		n.Url = Blog.Url
//		n.Sport = Blog.Sport
//		n.Featured = Blog.Featured
//		n.Exclusive = Blog.Exclusive
//		n.Trending = Blog.Trending
//		n.Content = Blog.Content
//		n.Picture = Blog.Picture
//		n.Code = Blog.Code
//		n.Base = Blog.Base
//		n.Comments = Blog.Comments
//		//sections
//		results := []Coder{}
//		for _, g := range Blog.Sections {
//			results = append(results, Coder{Name: g.Code})
//		}
//		n.Sections = results
//		return &n
//	}
// func (r blogrepository) GetBlogections(Blog *Blog) *Blog {
// 	// for _, m := range Blog.SectionsCodes {
// 	// 	res, _ := blogsections.BlogSectionrepository.GetOne(m)
// 	// 	Blog.Sections = append(Blog.Sections, res)
// 	// }
// 	r.ModifyOneBlogWithCategory(Blog)
// 	return Blog
// }

//	func (r blogrepository) EvaluateIfSecionExist(code string, oldinfo []*blogsections.BlogSection) (string, bool) {
//		var results bool = false
//		for _, g := range oldinfo {
//			if g.Code == code {
//				results = true
//				break
//			}
//		}
//		return code, results
//	}
//
//	func (r blogrepository) GetIfSecionExist(code string, newInfo []*blogsections.BlogSection) *blogsections.BlogSection {
//		var results blogsections.BlogSection
//		for _, g := range newInfo {
//			if g.Code == code {
//				results = *g
//				break
//			}
//		}
//		return &results
//	}
func (r blogrepository) DeleteAsection(code string, secs []Coder) []Coder {
	l := len(secs) - 1
	Bloglice := make([]Coder, l, l)
	for _, g := range secs {
		if g.Name != code {
			Bloglice = append(Bloglice, g)
		}
	}
	return Bloglice
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
