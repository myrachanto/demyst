package dashboard

import (
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/api/majorcategory"
	"github.com/myrachanto/estate/src/api/product"
	"github.com/myrachanto/estate/src/api/seo"
	subLocation "github.com/myrachanto/estate/src/api/sublocation"
)

type Dashboard struct {
	Products     []*product.Product `json:"products"`
	All          Modulo             `json:"all"`
	Featured     Modulo             `json:"featured"`
	Promoted     Modulo             `json:"trending"`
	HotDeals     Modulo             `json:"exclusive"`
	Chartdata    ChartData          `json:"chartdata"`
	Distribution []product.Modular  `json:"distribution"`
	Types        []product.Modular  `json:"types"`
	Linechart    []*product.Weekly  `json:"linechart"`
}

type Modulo struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}
type Module struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}
type Module2 struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Total int    `json:"total"`
}
type ChartData struct {
	Hotdeals Module `json:"exclusive"`
	Latest   Module `json:"latest"`
	Promoted Module `json:"trending"`
	Featured Module `json:"featured"`
	All      Module `json:"all"`
}
type Home struct {
	Featured   []product.Product   `json:"featured"`
	Promoted   []*product.Product  `json:"promoted"`
	Properties []Module2           `json:"properties"`
	Locations  []location.Location `json:"locations"`
	Seo        *seo.Seo            `json:"seo"`
}
type Nav struct {
	Locations     []location.Location           `json:"locations"`
	Sublocation   []subLocation.SubLocation     `json:"sublocations"`
	Majorcategory []majorcategory.Majorcategory `json:"majorcategories"`
}
