package product

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/api/majorcategory"
	"github.com/myrachanto/estate/src/api/seo"
	subLocation "github.com/myrachanto/estate/src/api/sublocation"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name"`
	Url           string             `json:"url"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Meta          string             `json:"meta"`
	Altertag      string             `json:"altertag"`
	Footer        string             `json:"footer"`
	Code          string             `json:"code"`
	Majorcategory string             `json:"majorcategory"`
	Category      string             `json:"category"`
	Oldprice      float64            `json:"oldprice"`
	Newprice      float64            `json:"newprice"`
	Buyprice      float64            `json:"buyprice"`
	Price         float64            `json:"price"`
	Likes         int64              `json:"likes"`
	Picture       string             `json:"picture"`
	Location      string             `json:"location"`
	SubLocation   string             `json:"sublocation"`
	Kind          string             `json:"kind"`
	Bathrooms     int64              `json:"bathrooms"`
	Bedrooms      int64              `json:"bedrooms"`
	Sqft          float64            `json:"sqft"`
	Length        float64            `json:"length"`
	Width         float64            `json:"width"`
	Video         string             `json:"video"`
	Services      []Service          `json:"services"`
	Rating        *Rating            `json:"rating"`
	Rates         []Rating           `json:"rates"`
	Images        []Picture          `json:"images"`
	Tag           []Tag              `json:"tags"`
	Features      []Feature          `json:"features"`
	Featured      bool               `json:"featured"`
	Promotion     bool               `json:"promotion"`
	Hotdeals      bool               `json:"hotdeals"`
	Complete      bool               `json:"complete"`
	Sold          bool               `json:"sold"`
	Base          support.Base       `json:"base,omitempty"`
}
type Results2 struct {
	Products []Product `json:"results"`
	Seo      seo.Seo   `json:"seo"`
}
type Rental struct {
	Product  []Product         `json:"product"`
	Seo      seo.Seo           `json:"seo"`
	Location location.Location `json:"location"`
}
type Response struct {
	Products []Product         `json:"products"`
	Seo      seo.Seo           `json:"seo"`
	Major    location.Location `json:"major"`
}
type Response2 struct {
	Products []Product               `json:"products"`
	Seo      seo.Seo                 `json:"seo"`
	Major    subLocation.SubLocation `json:"major"`
}
type Navs struct {
	Sublocations  []subLocation.SubLocation     `json:"sublocations"`
	Location      []location.Location           `json:"locations"`
	Majorcategory []majorcategory.Majorcategory `json:"majorcategories"`
}
type Property struct {
	Module []Module2 `json:"results"`
	Seo    seo.Seo   `json:"seo"`
}
type Module2 struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Total int    `json:"total"`
}

type ProductRelated struct {
	Product *Product   `json:"product"`
	Related []*Product `json:"related"`
}
type Rating struct {
	Author      string `json:"author"`
	Bestrate    int64  `json:"bestrate"`
	Rate        int64  `json:"rate"`
	Description string `json:"description"`
	TotalCount  int    `json:"totalcount"`
}
type Producto struct {
	Product     *Product  `json:"product"`
	Related     []Product `json:"related"`
	Bylocations []Product `json:"bylocation"`
}
type Newarrivals struct {
	Product []*Product `json:"product"`
}
type Service struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type Majorcat struct {
	Name string `json:"name"`
}
type Supercat struct {
	Name string `json:"name"`
}
type Categs struct {
	Name string `json:"name"`
}
type Tag struct {
	Name string `json:"name" `
	Code string `json:"code"`
}

type Feature struct {
	Name string `json:"name" `
	Code string `json:"code"`
}
type Picture struct {
	Productcode string `json:"productcode"`
	Name        string `json:"name"`
}
type Results struct {
	Data        []Product `json:"results"`
	Total       int       `json:"total"`
	Pages       int       `json:"pages"`
	CurrentPage int       `json:"currentpage"`
}

type ProductCount struct {
	Sport string `bson:"_id" json:"sport" `
	Count int    `json:"count"`
}
type Weekly struct {
	Id    int32 `bson:"_id" json:"id"`
	Count int32 `bson:"count" json:"count"`
}
type Modular struct {
	Id    string `bson:"_id" json:"id"`
	Total int64  `bson:"total" json:"total"`
}

func (product Product) Validate() httperrors.HttpErr {
	if product.Name == "" {
		return httperrors.NewNotFoundError("Invalid Name")
	}
	if product.Title == "" {
		return httperrors.NewNotFoundError("Invalid title")
	}
	// if product.Description == "" {
	// 	return httperrors.NewNotFoundError("Invalid Description")
	// }
	return nil
}
