package product

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/api/feature"
	feat "github.com/myrachanto/estate/src/api/feature"
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/api/majorcategory"
	"github.com/myrachanto/estate/src/api/rating"
	"github.com/myrachanto/estate/src/api/seo"
	subLocation "github.com/myrachanto/estate/src/api/sublocation"
	tag "github.com/myrachanto/estate/src/api/tags"
	"github.com/myrachanto/estate/src/db"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// productrepository ...
var (
	Productrepository ProductRepoInterface = &productrepository{}
	ctx                                    = context.TODO()
	Productos                              = productrepository{}
)

type ProductRepoInterface interface {
	Create(product *Product) httperrors.HttpErr
	GetOneE(code string) (*Producto, httperrors.HttpErr)
	GetOneE1(searcher support.Paginator2) (*Response, httperrors.HttpErr)
	GetOneE2(searcher support.Paginator2) (*Response2, httperrors.HttpErr)
	GetLand(searcher support.Paginator2) (*Response, httperrors.HttpErr)
	GetRental(searcher support.Paginator2) (*Response, httperrors.HttpErr)
	GetProperty(searcher support.Paginator2) (*Response, httperrors.HttpErr)
	GetPropertyType(searcher support.Paginator2) (*Response, httperrors.HttpErr)
	Search(search support.Paginator2) (*Response, httperrors.HttpErr)
	Search2(search support.Paginator2) ([]Product, httperrors.HttpErr)

	GetProductsbyMajorcategory(paginator *support.Paginator) (*Results, httperrors.HttpErr)
	GetProductsbycategory(search support.Paginator2) (*Response, httperrors.HttpErr)
	GroupbyLocation() ([]Modular, httperrors.HttpErr)
	GroupbyType() ([]Modular, httperrors.HttpErr)

	GetProperties() (*Property, httperrors.HttpErr)
	GetProductsFlavours(code string) ([]*Product, httperrors.HttpErr)
	GetProductsbyarrival(code string) (*Newarrivals, httperrors.HttpErr)
	GetProductshotdeals() (*Newarrivals, httperrors.HttpErr)
	GetProductshotdeals1(biz string) ([]*Product, httperrors.HttpErr)
	GetOne(code string) (*Product, httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	GetAll1(search string) ([]*Product, httperrors.HttpErr)
	Results(search support.Paginator) (*Results, httperrors.HttpErr)
	GetThree() ([]*Product, httperrors.HttpErr)
	GetFeatured() ([]Product, httperrors.HttpErr)
	GetProductsByLocation(loaction string) ([]Product, httperrors.HttpErr)

	GetNavs() (navs *Navs, errors httperrors.HttpErr)

	Update(code string, product *Product) httperrors.HttpErr
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateSold(code string, status bool) httperrors.HttpErr
	UpdateCompleted(code string, status bool) httperrors.HttpErr
	UpdateHotdeals(code string, status bool) httperrors.HttpErr
	UpdatePromotion(code string, status bool) httperrors.HttpErr
	Likes(code string, likes int64) httperrors.HttpErr
	AUpdate(code string, b, old, new, buy float64) httperrors.HttpErr
	Delete(id string) (string, httperrors.HttpErr)
	Count() (float64, httperrors.HttpErr)
	Feature() ([]*Product, httperrors.HttpErr)
	Hotdeal() ([]*Product, httperrors.HttpErr)
	Promotions(int) ([]*Product, httperrors.HttpErr)
	// GetProductsarrival() ([]*Product, httperrors.HttpErr)
	Producto1(code string) (product *Product)
	Producto(code string) (product *Product)
	Getratings(code string) (*Rating, []Rating)
	GetAllCategory() ([]*ProductCount, httperrors.HttpErr)
	GetAllPostByWeek() ([]*Weekly, httperrors.HttpErr)
	GetProductsByTag(biz, tag string) ([]*Product, httperrors.HttpErr)
	GetSize(size string) ([]Product, httperrors.HttpErr)
}
type productrepository struct {
	Mongodb *mongo.Database

	Cancel context.CancelFunc
}

func NewProductRepo() ProductRepoInterface {
	return &productrepository{}
}

func (r *productrepository) Create(product *Product) httperrors.HttpErr {

	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>step1", product)
	if err1 := product.Validate(); err1 != nil {
		return err1
	}
	// fmt.Println("step1 ---------------------------")
	code, err1 := r.genecode()
	if err1 != nil {
		return err1
	}
	// fmt.Println("step2 ---------------------------")
	product.Base.Updated_At = time.Now()
	product.Base.Created_At = time.Now()
	product.Code = code
	tags := []Tag{}
	for _, v := range product.Tag {
		tg := Tag{}
		str := strings.Join(strings.Split(v.Name, "-"), " ")
		// fmt.P
		t, err := tag.Tagrepository.GetOneByName(str)
		if err != nil {
			fmt.Println("ffffff>>>>>>>>>>>>", err, v.Name)
		}
		tg.Code = t.Code
		tg.Name = v.Name
		tags = append(tags, tg)
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>fffffffffffffff", product.Features)
	product.Tag = tags
	fets := []Feature{}
	for _, v := range product.Features {
		fet := Feature{}
		// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>fffffffffffffstr", v.Name)
		str := strings.Join(strings.Split(v.Name, "-"), " ")
		// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>fffffffffffffstr", str)
		f, err := feat.Featurerepo.GetOneByName(str)
		if err != nil {
			fmt.Println("ffffff>>>>>>>>>>>>", err)
		}
		fet.Code = f.Code
		fet.Name = v.Name
		fets = append(fets, fet)
	}
	product.Features = fets
	major, errs := majorcategory.Major.GetOnebyName(product.Majorcategory)
	if errs != nil {
		return errs
	}
	product.Majorcategory = major.Code
	locat, errs := location.Locationrepo.GetOneByName(product.Location)
	if errs != nil {
		return errs
	}
	product.Location = locat.Code
	if product.SubLocation != "" {
		subloc, errs := subLocation.SubLocationrepo.GetOneByName(product.SubLocation)
		if errs != nil {
			return errs
		}
		product.SubLocation = subloc.Code
	}

	collection := db.Mongodb.Collection("product")
	_, err := collection.InsertOne(ctx, &product)

	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Create product Failed, %d", err))
	}
	return nil
}

func (r *productrepository) GetOneE(code string) (*Producto, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"url", code}},
		}},
	}
	var product Product
	err := collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	product = *r.Assist(&product)
	related, errs := r.GetFeatured()
	if errs != nil {
		return nil, errs
	}
	bylocation, _ := r.GetProductsByLocation(product.Location)
	bylocations := []Product{}
	for _, g := range bylocation {
		res := r.Assist(&g)
		bylocations = append(bylocations, *res)
	}
	return &Producto{
		Product:     &product,
		Related:     related,
		Bylocations: bylocations,
	}, nil
}
func (r *productrepository) Assist(product *Product) *Product {
	major, errs := majorcategory.Major.Getuno(product.Majorcategory)
	if errs != nil {
		return nil
	}
	product.Majorcategory = major.Name
	locat, errs := location.Locationrepo.Getuno(product.Location)
	if errs != nil {
		return nil
	}
	product.Location = locat.Name
	name := ""
	sublocat, errs := subLocation.SubLocationrepo.Getuno(product.SubLocation)
	if errs == nil {
		name = sublocat.Name
	}
	product.SubLocation = name
	return product
}
func (r *productrepository) GetOneE1(searcher support.Paginator2) (*Response, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(searcher.Kind)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}
	town, err1 := location.Locationrepo.GetOneByName(searcher.Kind)
	if err1 != nil {
		return nil, err1
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"location", town.Code}},
		}},
	}
	products := []Product{}
	options := options.Find()
	skipNum := (searcher.Page - 1) * searcher.Pagesize
	options.SetLimit(int64(searcher.Pagesize))
	options.SetSkip(int64(skipNum))
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}

	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}
	seo := seo.Seorepo.GetOneByName(town.Name)
	return &Response{
		Products: prods,
		Seo:      *seo,
	}, nil
}

func (r *productrepository) GetOneE2(searcher support.Paginator2) (*Response2, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(searcher.Kind)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}
	town, err1 := location.Locationrepo.GetOneByName(searcher.Location)
	if err1 != nil {
		return nil, err1
	}
	subloc, errs := subLocation.SubLocationrepo.GetOneByName(searcher.Sublocation)
	if errs != nil {
		return nil, errs
	}
	var filter primitive.D
	if searcher.Kind == "properties" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", town.Code}},
				bson.D{{"sublocation", subloc.Code}},
				bson.D{{"kind", bson.D{{"$ne", "land"}}}},
			}},
		}
	}
	if searcher.Kind == "land" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", town.Code}},
				bson.D{{"sublocation", subloc.Code}},
				bson.D{{"kind", "land"}},
			}},
		}
	}
	if searcher.Kind == "rent" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", town.Code}},
				bson.D{{"sublocation", subloc.Code}},
				bson.D{{"kind", "rent"}},
			}},
		}
	}

	products := []Product{}
	options := options.Find()
	skipNum := (searcher.Page - 1) * searcher.Pagesize
	options.SetLimit(int64(searcher.Pagesize))
	options.SetSkip(int64(skipNum))
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}

	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}
	if searcher.Kind == "properties" {
		searcher.Kind = "sale"
	}
	// fmt.Println("---------------------------------step1")
	seo := seo.Seorepo.GetOneByName(searcher.Sublocation, searcher.Kind, town.Code)

	// fmt.Println("---------------------------------step4", seo)
	return &Response2{
		Products: prods,
		Seo:      *seo,
		Major:    *subloc,
	}, nil
}

func (r *productrepository) GetLand(searcher support.Paginator2) (*Response, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(searcher.Kind)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}
	town, err1 := location.Locationrepo.GetOneByName(searcher.Kind)
	if err1 != nil {
		return nil, err1
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"location", town.Code}},
			bson.D{{"kind", "land"}},
		}},
	}
	products := []Product{}
	options := options.Find()
	skipNum := (searcher.Page - 1) * searcher.Pagesize
	options.SetLimit(int64(searcher.Pagesize))
	options.SetSkip(int64(skipNum))
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}

	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}

	seo := seo.Seorepo.GetOneByName(town.Name, "land", town.Code)
	return &Response{
		Products: prods,
		Seo:      *seo,
		Major:    *town,
	}, nil

}
func (r *productrepository) GetRental(searcher support.Paginator2) (*Response, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(searcher.Kind)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}
	town, err1 := location.Locationrepo.GetOneByName(searcher.Kind)
	if err1 != nil {
		return nil, err1
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"location", town.Code}},
			bson.D{{"kind", "rent"}},
		}},
	}
	products := []Product{}
	options := options.Find()
	skipNum := (searcher.Page - 1) * searcher.Pagesize
	options.SetLimit(int64(searcher.Pagesize))
	options.SetSkip(int64(skipNum))
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}

	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}
	seo := seo.Seorepo.GetOneByName(town.Name, "rent", town.Code)
	return &Response{
		Products: prods,
		Seo:      *seo,
		Major:    *town,
	}, nil
}
func (r *productrepository) GetProperty(searcher support.Paginator2) (*Response, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(searcher.Kind)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}
	town, err1 := location.Locationrepo.GetOneByName(searcher.Kind)
	if err1 != nil {
		return nil, err1
	}
	// fmt.Println("+++++++++++++++++++++>>>>>>>>>>>>>", town)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"location", town.Code}},
			bson.D{{"kind", bson.D{{"$ne", "land"}}}},
		}},
	}
	products := []Product{}
	options := options.Find()
	skipNum := (searcher.Page - 1) * searcher.Pagesize
	options.SetLimit(int64(searcher.Pagesize))
	options.SetSkip(int64(skipNum))
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}

	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}
	seo := seo.Seorepo.GetOneByName(town.Name, "sale", town.Code)
	return &Response{
		Products: prods,
		Seo:      *seo,
		Major:    *town,
	}, nil
}
func (r *productrepository) GetPropertyType(searcher support.Paginator2) (*Response, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(searcher.Kind)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"url": code}
	major, err1 := majorcategory.Major.GetOnebyName(searcher.Kind)
	if err1 != nil {
		return nil, err1
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"majorcategory", major.Code}},
		}},
	}
	// filter := bson.D{
	// 	{"$and", bson.A{
	// 		bson.D{{"majorcategory", trimmed}},
	// 	}},
	// }
	products := []Product{}
	options := options.Find()
	skipNum := (searcher.Page - 1) * searcher.Pagesize
	options.SetLimit(int64(searcher.Pagesize))
	options.SetSkip(int64(skipNum))
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}

	// fmt.Println("----------------------resl", len(products))
	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}
	seo := seo.Seorepo.GetOneByName(major.Name)
	return &Response{
		Products: prods,
		// Location:     *town,
		Seo: *seo,
	}, nil
}
func (r *productrepository) Producto(code string) (product *Product) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil
	}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return
	}
	product = r.Assist(product)
	product.Rating, product.Rates = r.Getratings(product.Code)

	return product
}
func (r *productrepository) Producto1(code string) (product *Product) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil
	}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return
	}
	product = r.Assist(product)
	product.Rating, product.Rates = r.Getratings(product.Code)
	return product
}
func (r *productrepository) GetProductsbyMajorcategory(p *support.Paginator) (*Results, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(p.Search)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>adsvdvbs code")
	products := []Product{}
	// totalproducts := []Product{}

	// Calculate the number of documents to skip
	skipNum := (p.Page - 1) * p.Pagesize
	collection := db.Mongodb.Collection("product")
	res := strings.Split(p.Search, "-")
	if res[0] == "all" || res[0] == "All" || res[0] == "ALL" {

		filter := bson.D{}
		products := []Product{}
		findOptions := options.Find()
		findOptions.SetLimit(int64(p.Pagesize))
		findOptions.SetSkip(int64(skipNum))
		findOptions.SetSort(bson.D{{"name", 1}})
		cursor, err := collection.Find(ctx, filter, findOptions)
		// fmt.Println(cursor)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>adsvdvbs step 1")
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>adsvdvbs step 2", products[0].Rating)
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}

		// prods := r.Paginate(products, p.Page, p.Pagesize, len(products))
		return &Results{
			Data:  products,
			Total: int(count),
			Pages: int(count) / p.Pagesize,
		}, nil

	} else {

		findOptions := options.Find()
		findOptions.SetLimit(int64(p.Pagesize))
		findOptions.SetSkip(int64(skipNum))
		findOptions.SetSort(bson.D{{"name", 1}})
		// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>adsvdvbs code //////////////////", p.Search)
		filter := bson.D{
			{"$and", bson.A{
				bson.M{"supercategory": primitive.Regex{Pattern: p.Search, Options: "i"}},
			}},
		}
		// total, _ := collection.CountDocuments(ctx, filter)
		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		// fmt.Println("..................step1 -----------", len(products))

		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		// prods := r.Paginate(products, p.Page, p.Pagesize, len(products))
		// firstcourse, _ := r.CheckTags(p.Search)
		// // fmt.Println("..................step2]]]]]]]]]]]]]]]]]]]]]", len(products))
		// if len(firstcourse) == 0 {
		// 	fmt.Println("..................step2", len(products))
		// 	return &Results{
		// 		Data:  prods,
		// 		Total: len(products),
		// 	}, nil
		// }
		// // fmt.Println("..................step3", len(firstcourse))
		// for _, b := range firstcourse {
		// 	ok := r.EvaluateIfProductExits(b, products)
		// 	if !ok {
		// 		totalproducts = append(totalproducts, *b)
		// 	} else {

		// 		// fmt.Println(b)
		// 		totalproducts = append(totalproducts, *b)
		// 	}
		// }
		// fmt.Println("..................step4", len(products), len(firstcourse), len(totalproducts))
		// prods = r.Paginate(totalproducts, p.Page, p.Pagesize, len(totalproducts))
		// fmt.Println("-----------mxiiiiiiiiiii>>>>>>>>>>>>>>>>>", len(prods))
		// return &Results{
		// 	Data:  prods,
		// 	Total: len(totalproducts),
		// 	Pages: int(len(products) / p.Pagesize),
		// }, nil
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		// prods := r.Paginate(products, p.Page, p.Pagesize, len(products))
		return &Results{
			Data:  products,
			Total: int(count),
			Pages: int(count) / p.Pagesize,
		}, nil

	}
}

func (r *productrepository) EvaluateIfProductExits(p *Product, data []Product) bool {
	res := false
	for _, v := range data {
		if p.Code == v.Code {
			res = true
		} else {
			res = false
			return res
		}
	}
	return res
}

func (r *productrepository) Getratings(code string) (*Rating, []Rating) {
	ratings, err := rating.Ratos.GetRatings(code)
	if err != nil {
		return nil, nil
	}
	var rating = &Rating{}
	var rates = Rating{}
	var ratess = []Rating{}
	var rato int64 = 0
	var bestra int64 = 0
	var count int
	for _, k := range ratings {
		count++
		rato += k.Rate
		if bestra < k.Rate {
			rating.Bestrate = k.Rate
			rating.Author = k.Author
			rating.Description = k.Description
		}
		rates.Author = k.Author
		rates.Description = k.Description
		rates.Rate = k.Rate
		ratess = append(ratess, rates)
	}
	rating.Rate = rato
	rating.TotalCount = count
	return rating, ratess
}

//	func (r *productrepository) Paginate(x []Product, skip int, size int) []Product {
//		if skip > len(x) {
//			skip = len(x)
//		}
//		end := skip + size
//		if end > len(x) {
//			end = len(x)
//		}
//		return x[skip:end]
//	}
func (r *productrepository) Paginate(x []Product, pageNum int, pageSize int, sliceLength int) []Product {
	start := 0
	if pageNum > 1 {
		s := pageNum - 1
		start = s * pageSize
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>Steps", pageNum, pageSize)
	if start > sliceLength {
		start = sliceLength
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>Steps", start, pageNum, pageSize)

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return x[start:end]
}

// func (r *productrepository) CheckUniqueness(data []Product) []Product {
// 	results := []Product{}
// 	for _, _ = range data {
// 		if len(results) == 0 {
// 			results = append(results, data[0])
// 			fmt.Println("..................step1", len(results))
// 		} else {
// 			for _, g := range results {
// 				for _, k := range data {
// 					if g.Code == k.Code {
// 						fmt.Println("..................step2", len(results))
// 						//do nothing
// 						break
// 					} else {
// 						fmt.Println("..................step3", len(results))
// 						results = append(results, k)
// 					}
// 				}

// 				fmt.Println("..................step3", len(results))
// 			}

//			}
//		}
//		return results
//	}
func (r *productrepository) CheckTags(name string) (product []*Product, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(name)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	res := strings.Split(name, "-")
	results := []*Product{}
	products := []*Product{}
	for _, g := range res {
		filter := bson.D{
			{"$and", bson.A{
				bson.M{"tag.name": primitive.Regex{Pattern: g, Options: "i"}},
			}},
		}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		products = append(products, results...)

	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println("<<<<<<<<<<<<<<<stteepeepee2", products)
	return products, nil

}
func (r *productrepository) GetProperties() (*Property, httperrors.HttpErr) {
	var property Property

	properties, err := r.GroupbyMajorcategory()
	if err != nil {
		return nil, err
	}
	mods := []Module2{}
	mod := Module2{}
	for _, v := range properties {
		major, _ := majorcategory.Major.GetOne(v.Id)
		mod.Name = major.Name
		mod.Image = major.Picture
		mod.Total = int(v.Total)
		mods = append(mods, mod)
	}
	property.Module = mods

	fmt.Println("++++++++++++++++++++++step 1 propasasd ----------")
	seo := seo.Seorepo.GetOneByName("properties")

	property.Seo = *seo
	return &property, nil

}
func (r *productrepository) GetProductsbycategory(search support.Paginator2) (*Response, httperrors.HttpErr) {
	fmt.Println("++++++++++++++++++++++step 1")
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"category": code}
	code := search.Kind
	if code == "rental" {
		code = "rent"
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"kind", code}},
		}},
	}
	if code == "properties" {
		// code = "land"
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"kind", bson.D{{"$ne", "land"}}}},
			}},
		}

	}
	products := []*Product{}
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, opts)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	prods := []Product{}
	for _, pr := range products {
		p := r.Assist(pr)
		prods = append(prods, *p)
	}
	// fmt.Println("+++++++++++++++++++++1", code)
	if code == "properties" {
		code = "sale"
	}
	// fmt.Println("+++++++++++++++++++++2", code)
	seo := seo.Seorepo.GetOneByName(code)
	return &Response{
		Products: prods,
		Seo:      *seo,
	}, nil
}

func (r *productrepository) GetProductsFlavours(url string) ([]*Product, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(url)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	// fmt.Println("++++++++++++++++++++++++++", url)
	tag, errs := tag.Tagrepository.GetOneByUrl(url)
	if errs != nil {
		return nil, errs
	}

	// fmt.Println("++++++++++++++++++++++++++code", tag.Code)
	collection := db.Mongodb.Collection("product")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"tag.code", tag.Code}},
		}},
	}
	// filter := bson.D{
	// 	{Key: "supercategory", Value: bson.D{{Key: "name", Value: code}}},
	// }
	// filter := bson.D{
	// 	{"$and", bson.A{
	// 		bson.D{{"shopalias", }},
	// 	}},
	// 	{Key: "supercategory", Value: bson.D{{Key: "name", Value: code}}},
	// }
	// pipeline := []bson.M{
	// 	{
	// 		"$match": bson.M{
	// 			"$and": []bson.M{
	// 				{
	// 					"shopalias": ,
	// 				},
	// 				{
	// 					"items": bson.M{
	// 						"$elemMatch": bson.M{
	// 							"code": code,
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	products := []*Product{}
	cursor, err := collection.Find(ctx, filter)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return products, nil
}
func (r *productrepository) GetProductsByTag(biz, tag string) ([]*Product, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(tag)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection("product")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"tag.code", tag}},
		}},
	}

	products := []*Product{}
	cursor, err := collection.Find(ctx, filter)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return products, nil
}
func (r *productrepository) GroupbyMajorcategory() ([]Modular, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")

	// Define the group stage
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$majorcategory"},
			{"total", bson.D{
				{"$sum", 1},
			}},
		}},
	}

	// Perform aggregation
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}

	var results []Modular
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println("sdsf========> ", err)
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>group major category", results)

	return results, nil
}

func (r *productrepository) GroupbyLocation() ([]Modular, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")

	// Define the group stage
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$location"},
			{"total", bson.D{
				{"$sum", 1},
			}},
		}},
	}

	// Perform aggregation
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}

	var results []Modular
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println("sdsf========> ", err)
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>group major category", results)

	return results, nil
}
func (r *productrepository) GroupbyType() ([]Modular, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")

	// Define the group stage
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$kind"},
			{"total", bson.D{
				{"$sum", 1},
			}},
		}},
	}

	// Perform aggregation
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}

	var results []Modular
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println("sdsf========> ", err)
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>group major category", results)

	return results, nil
}
func (r *productrepository) GetProductsarrival(string) ([]*Product, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"newarrivals": "true"}
	// fmt.Println("--------------------------")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"newarrivals", "true"}},
		}},
	}
	products := []*Product{}
	options := options.Find()
	options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return products, nil
}
func (r *productrepository) Feature() ([]*Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"featured": true}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"featured", true}},
		}},
	}
	products := []*Product{}
	options := options.Find()
	// options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return products, nil

}
func (r *productrepository) Hotdeal() ([]*Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"hotdeals": true}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"hotdeals", true}},
		}},
	}
	products := []*Product{}
	options := options.Find()
	// options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return products, nil

}

func (r *productrepository) Promotions(limit int) ([]*Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"promotion": true}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"promotion", true}},
		}},
	}
	products := []*Product{}
	options := options.Find()
	if limit > 0 {
		options.SetLimit(2)
	}
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	prods := []*Product{}
	for _, p := range products {
		prd := r.Assist(p)
		prods = append(prods, prd)
	}
	return prods, nil

}
func (r *productrepository) GetProductsbyarrival(code string) (*Newarrivals, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}

	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"newarrivals": "true"}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"newarrivals", "true"}},
			bson.D{{"code", code}},
		}},
	}
	products := []*Product{}
	options := options.Find()
	options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return &Newarrivals{
		Product: products,
	}, nil
}
func (r *productrepository) GetProductshotdeals() (*Newarrivals, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"hotdeals": true}
	// fmt.Println(">>>>>>>>>>>>>>>>wow", )
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"hotdeals", true}},
		}},
	}

	products := []*Product{}
	options := options.Find()
	// options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return &Newarrivals{
		Product: products,
	}, nil
}
func (r *productrepository) GetProductshotdeals1(biz string) ([]*Product, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"hotdeals": true}
	// fmt.Println(">>>>>>>>>>>>>>>>wow", )
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"hotdeals", true}},
		}},
	}

	products := []*Product{}
	options := options.Find()
	// options.SetLimit(5)
	options.SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, filter, options)
	// fmt.Println(cursor)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	// fmt.Println(products)
	return products, nil
}

func (r *productrepository) GetNavs() (navs *Navs, errors httperrors.HttpErr) {

	locations, err := location.Locationrepo.GetAll()
	if err != nil {
		return nil, err
	}
	subLocaions, err := subLocation.SubLocationrepo.GetAll()
	if err != nil {
		return nil, err
	}
	majorcategories, err := majorcategory.Major.GetAll()
	return &Navs{
		Sublocations:  subLocaions,
		Location:      locations,
		Majorcategory: majorcategories,
	}, nil
}
func (r *productrepository) GetOne(code string) (product *Product, errors httperrors.HttpErr) {
	fmt.Println("---------------------->>>>>>>", code)
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}

	collection := db.Mongodb.Collection("product")
	// filter := bson.M{"code": code}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	fmt.Println("---------------------->>>>>>>2", product)
	product = r.Assist(product)
	fmt.Println("---------------------->>>>>>>3", product)
	return product, nil
}

func (r *productrepository) GetAll1(search string) ([]*Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []*Product{}
	fmt.Println(search)
	if search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search, Options: "i"}}},
				bson.D{{"title", primitive.Regex{Pattern: search, Options: "i"}}},
				bson.D{{"description", primitive.Regex{Pattern: search, Options: "i"}}},
			}},
		}
		// fmt.Println(filter)
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		return products, nil
	} else {
		filter := bson.D{}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return []*Product{}, nil
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		return products, nil
	}
}

func (r *productrepository) Search(search support.Paginator2) (*Response, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []Product{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	// fmt.Printf("step 1++++++++++++++++++++>>>>>>>>>>location: %s sublocation: %s and kind: %s \n", search.Location, search.Sublocation, search.Kind)

	// if search.Kind == "properties" {
	// 	search.Kind = "land"
	// 	filter = bson.D{
	// 		{"$and", bson.A{
	// 			bson.D{{"kind", bson.D{{"$ne", code}}}},
	// 		}},
	// 	}

	// }
	local := ""
	location, _ := location.Locationrepo.GetOneByName(search.Location)
	if location != nil {
		local = location.Code
	}
	sublocal := ""
	sublocation, _ := subLocation.SubLocationrepo.GetOneByName(search.Sublocation)
	if sublocation != nil {
		sublocal = sublocation.Code
	}
	major := ""
	majord, _ := majorcategory.Major.GetOnebyName(search.Kind)
	if sublocation != nil {
		major = majord.Code
	}
	var filter primitive.D
	if search.Location != "" && search.Sublocation != "" && search.Majorcat != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", local}},
				bson.D{{"sublocation", sublocal}},
				bson.D{{"majorcategory", major}},
			}},
		}
	}
	if search.Location != "" && search.Sublocation != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", local}},
				bson.D{{"sublocation", sublocal}},
			}},
		}
	}
	if search.Location != "" && search.Majorcat != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", local}},
				bson.D{{"majorcategory", major}},
			}},
		}
	}
	if search.Sublocation != "" && search.Majorcat != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"sublocation", sublocal}},
				bson.D{{"majorcategory", major}},
			}},
		}
	}
	if search.Location != "" && search.Sublocation == "" && search.Majorcat == "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", local}},
			}},
		}
	}
	if search.Location == "" && search.Sublocation == "" && search.Majorcat != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"majorcategory", major}},
			}},
		}
	}
	if search.Location == "" && search.Majorcat == "" && search.Sublocation != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"sublocation", sublocal}},
			}},
		}
	}
	if search.Location == "" && search.Sublocation == "" && search.Majorcat == "" {
		filter = bson.D{}
	}
	// Additional filter based on the condition
	if search.Kind == "rental" {
		search.Kind = "rent"
	}
	additionalFilter := bson.D{}
	if search.Kind == "properties" {
		search.Kind = "land"
		additionalFilter = bson.D{
			{"$and", bson.A{
				bson.D{{"kind", bson.D{{"$ne", search.Kind}}}},
			}},
		}
	}

	// Combine filters using $or operator
	filter2 := bson.D{
		{"$and", bson.A{filter, additionalFilter}},
	}

	// fmt.Printf("++++++++++++++++++++>>>>>>>>>>location: %s sublocation: %s and kind: %s \n", local, sublocal, major)
	cursor, err := collection.Find(ctx, filter2, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(&f)
		prods = append(prods, *prod)
	}
	// fmt.Println("++++++++++++++++++++>>>>>>>>>>prods", prods)
	return &Response{
		Products: prods,
	}, nil
}

func (r *productrepository) Search2(search support.Paginator2) ([]Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []Product{}
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	local := ""
	location, _ := location.Locationrepo.GetOneByName(search.Location)
	if location != nil {
		local = location.Code
	}
	major := ""
	majord, _ := majorcategory.Major.GetOnebyName(search.Majorcat)
	if majord != nil {
		major = majord.Code
	}
	var filter primitive.D
	if search.Location != "" && search.Majorcat != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", local}},
				bson.D{{"majorcategory", major}},
			}},
		}
	}
	if search.Location != "" && search.Majorcat == "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"location", local}},
			}},
		}
	}
	if search.Location == "" && search.Majorcat != "" {
		filter = bson.D{
			{"$and", bson.A{
				bson.D{{"majorcategory", major}},
			}},
		}
	}
	if search.Location == "" && search.Majorcat == "" {
		filter = bson.D{}
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	return products, nil
}
func (r *productrepository) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []Product{}
	// fmt.Println("---------------------->>>>>>>", search)
	// Calculate the number of documents to skip
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	if search.Search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search.Search, Options: "i"}}},
				bson.D{{"title", primitive.Regex{Pattern: search.Search, Options: "i"}}},
				bson.D{{"description", primitive.Regex{Pattern: search.Search, Options: "i"}}},
			}},
		}
		// fmt.Println(filter)
		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}

		prods := []Product{}
		for _, f := range products {
			prod := r.Assist(&f)
			prods = append(prods, *prod)
		}
		return &Results{
			Data:  prods,
			Total: int(count),
			Pages: int(count / int64(search.Pagesize)),
		}, nil
	} else {
		filter := bson.D{}
		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		prods := []Product{}
		for _, f := range products {
			prod := r.Assist(&f)
			prods = append(prods, *prod)
		}
		return &Results{
			Data:  prods,
			Total: int(count),
			Pages: int(count / int64(search.Pagesize)),
		}, nil
	}
}

func (r *productrepository) Results(search support.Paginator) (*Results, httperrors.HttpErr) {

	fmt.Println(">>>>>>>>>>>>>>>>>>>>step1", search)
	collection := db.Mongodb.Collection("product")
	products := []Product{}
	// Calculate the number of documents to skip
	skipNum := (search.Page - 1) * search.Pagesize
	findOptions := options.Find()
	findOptions.SetLimit(int64(search.Pagesize))
	findOptions.SetSkip(int64(skipNum))
	// fmt.Println(search)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>step2")
	if search.Search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"name", primitive.Regex{Pattern: search.Search, Options: "i"}}},
				bson.D{{"title", primitive.Regex{Pattern: search.Search, Options: "i"}}},
				bson.D{{"description", primitive.Regex{Pattern: search.Search, Options: "i"}}},
			}},
		}
		// fmt.Println(filter)
		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		return &Results{
			Data:        products,
			Total:       len(products),
			CurrentPage: search.Page,
		}, nil
	} else {
		filter := bson.D{
			{"$and", bson.A{
				bson.D{{search.Style, true}}},
			},
		}
		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &products); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		for _, pr := range products {
			pr.Rating, pr.Rates = r.Getratings(pr.Code)
		}
		fmt.Println(">>>>>>>>>>>>>>>>>>>>", products)
		return &Results{
			Data:        products,
			Total:       len(products),
			CurrentPage: search.Page,
		}, nil
	}

}
func (r *productrepository) GetThree() ([]*Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []*Product{}
	options := options.Find()

	// Limit by 10 documents only
	options.SetLimit(4)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"hotdeals", true}},
		}},
	}
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	return products, nil

}
func (r *productrepository) GetSize1(size string) ([]Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []Product{}
	options := options.Find()

	// Limit by 10 documents only
	options.SetLimit(4)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"services.name", size}},
		}},
	}
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	return products, nil

}
func (r *productrepository) GetSize(size string) ([]Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []Product{}
	options := options.Find()

	// Limit by 10 documents only
	// options.SetLimit(limit)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"services.name", size}},
		}},
	}
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	prods := []Product{}
	for _, f := range products {
		prod := Product{}
		for _, k := range f.Services {
			if size == k.Name {
				prod = f
				prod.Code = k.Code
				prods = append(prods, prod)
			}
		}

	}
	return prods, nil

}
func (r *productrepository) GetFeatured() ([]Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []*Product{}
	options := options.Find()

	// Limit by 10 documents only
	options.SetLimit(4)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"featured", true}},
		}},
	}
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	prods := []Product{}
	for _, f := range products {
		prod := r.Assist(f)
		prods = append(prods, *prod)
	}
	return prods, nil

}
func (r *productrepository) GetProductsByLocation(location string) ([]Product, httperrors.HttpErr) {

	collection := db.Mongodb.Collection("product")
	products := []Product{}
	options := options.Find()

	// Limit by 10 documents only
	// options.SetLimit(8)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"location", location}},
		}},
	}
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, httperrors.NewNotFoundError("No records found!")
	}
	if err = cursor.All(ctx, &products); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}
	for _, pr := range products {
		pr.Rating, pr.Rates = r.Getratings(pr.Code)
	}
	return products, nil

}
func (r *productrepository) GetAllCategory() ([]*ProductCount, httperrors.HttpErr) {
	results := []*ProductCount{}
	collection := db.Mongodb.Collection("product")
	fmt.Println("_____________________step 1")

	groupStage2 := bson.M{
		"$group": bson.M{
			"_id": "$majorcategory",
			"count": bson.M{
				"$sum": 1,
			},
		},
	}

	pipeline := []bson.M{groupStage2}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("_____________________step 2", err)
		return nil, httperrors.NewNotFoundError("No records found!")
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, httperrors.NewNotFoundError("Error decoding!")
	}

	return results, nil

}

func (r *productrepository) GetAllPostByWeek() ([]*Weekly, httperrors.HttpErr) {
	results := []*Weekly{}
	collection := db.Mongodb.Collection("product")
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
func (r *productrepository) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Product{}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
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
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) UpdateSold(code string, status bool) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Product{}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"sold", status}}},
		},
	)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) UpdateCompleted(code string, status bool) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Product{}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"complete", status}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) UpdateHotdeals(code string, status bool) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Product{}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"hotdeals", status}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) UpdatePromotion(code string, status bool) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Product{}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	// fmt.Println("------------------update promotions", code, status)
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"promotion", status}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) Likes(code string, likes int64) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	upay := &Product{}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"likes", likes}}},
		},
	)
	// update := bson.M{"$set": pay}
	// _, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) Update(code string, product *Product) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}
	fmt.Println("step 1 ============================", code)
	if product.Majorcategory != "" {
		major, errs := majorcategory.Major.GetOnebyName(product.Majorcategory)
		if errs != nil {
			return errs
		}
		product.Majorcategory = major.Code
	}
	fmt.Println("step 2 ============================")
	if product.Location != "" {
		locat, errs := location.Locationrepo.GetOneByName(product.Location)
		if errs != nil {
			return errs
		}
		product.Location = locat.Code
	}
	// fmt.Println("step 3 ============================")
	if product.SubLocation != "" {
		subloc, errs := subLocation.SubLocationrepo.GetOneByName(product.SubLocation)
		if errs != nil {
			return errs
		}
		product.SubLocation = subloc.Code
	}
	// fmt.Println("step 4 ============================")
	tags := []Tag{}
	tg := Tag{}
	// fmt.Println("------------------step1 tagiing",  product.Tag)
	if len(product.Tag) > 0 {
		for _, v := range product.Tag {
			str := strings.Join(strings.Split(v.Name, "-"), " ")
			t, _ := tag.Tagrepository.GetOneByName(str)
			tg.Code = t.Code
			tg.Name = v.Name
			tags = append(tags, tg)
		}

	}
	fets := []Feature{}
	fet := Feature{}
	if len(product.Features) > 0 {
		for _, v := range product.Features {
			str := strings.Join(strings.Split(v.Name, "-"), " ")
			t, _ := feature.Featurerepo.GetOneByName(str)
			fet.Code = t.Code
			fet.Name = v.Name
			fets = append(fets, fet)
		}

	}
	product.Features = fets
	product.Tag = tags
	// fmt.Println("------------------step1 tags ", tags)
	result, errs := r.getuno(code)
	if errs != nil {
		return errs
	}
	if len(product.Images) > 0 {
		for _, image := range result.Images {
			go support.Clean.Cleaner(image.Name)
		}
	}
	if product.Name == "" {
		product.Name = result.Name
	}
	if product.Title == "" {
		product.Title = result.Title
	}
	if product.Description == "" {
		product.Description = result.Description
	}
	if product.Code == "" {
		product.Code = result.Code
	}
	if product.Url == "" {
		product.Url = result.Url
	}
	if product.Meta == "" {
		product.Meta = result.Meta
	}
	if product.Altertag == "" {
		product.Altertag = result.Altertag
	}
	if product.Footer == "" {
		product.Footer = result.Footer
	}
	if !product.Featured {
		product.Featured = result.Featured
	}
	if !product.Promotion {
		product.Promotion = result.Promotion
	}
	if !product.Hotdeals {
		product.Hotdeals = result.Hotdeals
	}

	if product.Majorcategory == "" {
		product.Majorcategory = result.Majorcategory
	}
	if product.Category == "" {
		product.Category = result.Category
	}
	if product.Oldprice == 0 {
		product.Oldprice = result.Oldprice
	}
	if product.Newprice == 0 {
		product.Newprice = result.Newprice
	}
	if product.Buyprice == 0 {
		product.Buyprice = result.Buyprice
	}
	if product.Price == 0 {
		product.Price = result.Price
	}
	if product.Sqft == 0 {
		product.Sqft = result.Sqft
	}
	if product.Bathrooms == 0 {
		product.Bathrooms = result.Bathrooms
	}
	if product.Bedrooms == 0 {
		product.Bedrooms = result.Bedrooms
	}
	if product.Buyprice == 0 {
		product.Buyprice = result.Buyprice
	}
	if product.Width == 0 {
		product.Width = result.Width
	}
	if product.Length == 0 {
		product.Length = result.Length
	}
	if product.Picture == "" {
		product.Picture = result.Picture
	}
	if product.Video == "" {
		product.Video = result.Video
	}
	if product.Kind == "" {
		product.Kind = result.Kind
	}
	if product.SubLocation == "" {
		product.SubLocation = result.SubLocation
	}

	if product.Location == "" {
		product.Location = result.Location
	}

	if len(product.Images) == 0 {
		product.Images = result.Images
	}
	if len(product.Services) == 0 {
		product.Services = result.Services
	}
	if len(product.Tag) == 0 {
		product.Tag = result.Tag
	}
	if len(product.Features) == 0 {
		product.Features = result.Features
	}
	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	product.Base.Updated_At = time.Now()
	update := bson.M{"$set": product}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}

func (r *productrepository) AUpdate(code string, b, old, new, buy float64) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return stringresults
	}

	result, err3 := r.getuno(code)
	if err3 != nil {
		fmt.Println(err3)
	}
	product := &Product{}
	product.Oldprice = old
	product.Newprice = new
	product.Buyprice = buy
	product.Base.Updated_At = time.Now()
	if product.Name == "" {
		product.Name = result.Name
	}
	if product.Title == "" {
		product.Title = result.Title
	}
	if product.Description == "" {
		product.Description = result.Description
	}
	if product.Code == "" {
		product.Code = result.Code
	}

	if product.Oldprice == 0 {
		product.Oldprice = result.Oldprice
	}
	if product.Newprice == 0 {
		product.Newprice = result.Newprice
	}
	if product.Buyprice == 0 {
		product.Buyprice = result.Buyprice
	}

	collection := db.Mongodb.Collection("product")
	filter := bson.M{"code": code}
	update := bson.M{"$set": product}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return httperrors.NewNotFoundError("Error updating!")
	}
	return nil
}
func (r *productrepository) Delete(id string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return "", stringresults
	}
	prod, errs := r.getuno(id)
	if errs != nil {
		return "", errs
	}
	for _, image := range prod.Images {
		go support.Clean.Cleaner(image.Name)
	}

	collection := db.Mongodb.Collection("product")

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", id}},
		}},
	}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil
}
func (r *productrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp
	collection := db.Mongodb.Collection("product")
	filter := bson.D{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	code := "ProductCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special

	return code, nil
}
func (r *productrepository) getuno(code string) (result *Product, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}

	collection := db.Mongodb.Collection("product")
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"code", code}},
		}},
	}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r *productrepository) Count() (float64, httperrors.HttpErr) {
	collection := db.Mongodb.Collection("product")
	filter := bson.D{}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil
	}
	code := float64(count)

	return code, nil
}
