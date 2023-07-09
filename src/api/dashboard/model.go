package dashboard

import (
	"github.com/myrachanto/sports/src/api/news"
	"github.com/myrachanto/sports/src/api/pages"
)

type Dashboard struct {
	News        []*news.News       `json:"news"`
	All         Module             `json:"all"`
	Trending    Module             `json:"trending"`
	Exclusive   Module             `json:"exclusive"`
	Featured    Module             `json:"featured"`
	Chartdata   ChartData          `json:"chartdata"`
	Sportcounts []*news.SportCount `json:"sportcount"`
	Linechart   []*news.Weekly     `json:"linechart"`
}
type Module struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}
type ChartData struct {
	Trending  Module `json:"trending"`
	Latest    Module `json:"latest"`
	Exclusive Module `json:"exclusive"`
	Featured  Module `json:"featured"`
	All       Module `json:"all"`
}
type Home struct {
	All       []*news.News `json:"all,omitempty"`
	Latest    []*news.News `json:"latest,omitempty"`
	Trending  []*news.News `json:"trending,omitempty"`
	Exclusive []*news.News `json:"exclusive,omitempty"`
	Featured  []*news.News `json:"featured,omitempty"`
	Seo       *pages.Page  `json:"seo,omitempty"`
}
