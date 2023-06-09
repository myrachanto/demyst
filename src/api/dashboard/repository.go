package dashboard

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/api/news"
	"github.com/myrachanto/sports/src/api/pages"
)

// dashboardrepository repository
var (
	Dashboardrepository DashboardrepoInterface = &dashboardrepository{}
	Dashboardrepo                              = dashboardrepository{}
)

type DashboardrepoInterface interface {
	HomeCms() (*Dashboard, httperrors.HttpErr)
	Index() (*Home, httperrors.HttpErr)
}
type dashboardrepository struct{}

func NewdashboardRepo() DashboardrepoInterface {
	return &dashboardrepository{}
}

func (r *dashboardrepository) HomeCms() (*Dashboard, httperrors.HttpErr) {
	allnews, err := news.Newsrepo.GetAll1()
	if err != nil {
		return nil, err
	}
	trending, err := news.Newsrepo.GetAllTrending()
	if err != nil {
		return nil, err
	}
	exclusive, err := news.Newsrepo.GetAllExclusive()
	if err != nil {
		return nil, err
	}
	featured, err := news.Newsrepo.GetAllFeatured()
	if err != nil {
		return nil, err
	}
	sports, err := news.Newsrepo.GetAllCategory()
	if err != nil {
		return nil, err
	}
	weekly, err := news.Newsrepo.GetAllPostByWeek()
	if err != nil {
		return nil, err
	}
	var dash Dashboard
	dash.News = trending
	dash.All.Name = "All"
	dash.All.Total = len(allnews)
	dash.Trending.Name = "Trending"
	dash.Trending.Total = len(trending)
	dash.Exclusive.Name = "Exclusive"
	dash.Exclusive.Total = len(exclusive)
	dash.Featured.Name = "Featured"
	dash.Featured.Total = len(featured)
	dash.Chartdata.Trending.Name = "Trending"
	dash.Chartdata.Trending.Total = len(trending)
	dash.Chartdata.All.Name = "All"
	dash.Chartdata.All.Total = len(allnews)
	dash.Chartdata.Exclusive.Name = "Exclusive"
	dash.Chartdata.Exclusive.Total = len(exclusive)
	dash.Chartdata.Featured.Name = "Featured"
	dash.Chartdata.Featured.Total = len(featured)
	dash.Sportcounts = sports
	dash.Linechart = weekly

	return &dash, nil
}
func (r *dashboardrepository) Index() (*Home, httperrors.HttpErr) {
	allnews, err := news.Newsrepo.GetAll1()
	if err != nil {
		return nil, err
	}
	trending, err := news.Newsrepo.GetAllTrending(4)
	if err != nil {
		return nil, err
	}
	exclusive, err := news.Newsrepo.GetAllExclusive(4)
	if err != nil {
		return nil, err
	}
	featured, err := news.Newsrepo.GetAllFeatured(7)
	if err != nil {
		return nil, err
	}
	seo, err := pages.Pagerepo.GetOneByName("home")
	if err != nil {
		return nil, err
	}

	return &Home{
		Latest:    allnews,
		Featured:  featured,
		Exclusive: exclusive,
		Trending:  trending,
		Seo:       seo,
	}, nil
}
