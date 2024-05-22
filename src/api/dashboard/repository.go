package dashboard

import (
	"fmt"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/api/majorcategory"
	"github.com/myrachanto/estate/src/api/product"
	"github.com/myrachanto/estate/src/api/seo"
	subLocation "github.com/myrachanto/estate/src/api/sublocation"
)

// dashboardrepository repository
var (
	Dashboardrepository DashboardrepoInterface = &dashboardrepository{}
	Dashboardrepo                              = dashboardrepository{}
)

type DashboardrepoInterface interface {
	HomeCms() (*Dashboard, httperrors.HttpErr)
	Index() (*Home, httperrors.HttpErr)
	Layout() (*Nav, httperrors.HttpErr)
}
type dashboardrepository struct{}

func NewdashboardRepo() DashboardrepoInterface {
	return &dashboardrepository{}
}

func (r *dashboardrepository) HomeCms() (*Dashboard, httperrors.HttpErr) {
	// fmt.Println("----------------------------step 1")
	products, err := product.Productrepository.GetAll1("")
	if err != nil {
		return nil, err
	}
	// fmt.Println("----------------------------step 2")
	feature, err := product.Productrepository.Feature()
	if err != nil {
		return nil, err
	}
	// fmt.Println("----------------------------step 3")
	promotions, err := product.Productrepository.Promotions(0)
	if err != nil {
		return nil, err
	}
	// fmt.Println("----------------------------step 4")
	hotdeals, err := product.Productrepository.Hotdeal()
	if err != nil {
		return nil, err
	}
	// fmt.Println("----------------------------step 6")
	weekly, err := product.Productrepository.GetAllPostByWeek()
	if err != nil {
		return nil, err
	}
	locs, err := product.Productrepository.GroupbyLocation()
	if err != nil {
		return nil, err
	}
	mods := []product.Modular{}
	for _, l := range locs {
		loc, _ := location.Locationrepo.Getuno(l.Id)
		l.Id = loc.Name
		mods = append(mods, l)
	}

	types, err := product.Productrepository.GroupbyType()
	if err != nil {
		return nil, err
	}
	fmt.Println("----------------------------step 7", locs)
	var dash Dashboard
	dash.Products = feature
	dash.All.Name = "All"
	dash.All.Total = len(products)
	dash.Featured.Name = "Featured"
	dash.Featured.Total = len(feature)
	dash.Promoted.Name = "Promoted"
	dash.Promoted.Total = len(promotions)
	dash.HotDeals.Name = "HotDeals"
	dash.HotDeals.Total = len(hotdeals)
	dash.Chartdata.Featured.Name = "Featured"
	dash.Chartdata.Featured.Total = len(feature)
	dash.Chartdata.All.Name = "All"
	dash.Chartdata.All.Total = len(products)
	dash.Chartdata.Hotdeals.Name = "Hotdeals"
	dash.Chartdata.Hotdeals.Total = len(hotdeals)
	dash.Chartdata.Promoted.Name = "Promoted"
	dash.Chartdata.Promoted.Total = len(promotions)
	dash.Distribution = mods
	dash.Types = types
	dash.Linechart = weekly

	return &dash, nil
}
func (r *dashboardrepository) Index() (*Home, httperrors.HttpErr) {

	featured, err := product.Productos.GetFeatured()
	if err != nil {
		return nil, err
	}
	locations, err := location.Locationrepo.GetAll()
	if err != nil {
		return nil, err
	}
	seo := seo.Seorepo.GetOneByName("home")
	properties, err := product.Productos.GroupbyMajorcategory()
	if err != nil {
		return nil, err
	}

	// fmt.Println("----------------------------step 3")
	promotions, err := product.Productrepository.Promotions(2)
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

	return &Home{
		Featured:   featured,
		Promoted:   promotions,
		Locations:  locations,
		Properties: mods,
		Seo:        seo,
	}, nil
}

func (r *dashboardrepository) Layout() (*Nav, httperrors.HttpErr) {

	locations, err := location.Locationrepo.GetAll()
	if err != nil {
		return nil, err
	}
	sublocaions, err := subLocation.SubLocationrepo.GetAll()
	if err != nil {
		return nil, err
	}
	majorcategory, err := majorcategory.Major.GetAll()
	if err != nil {
		return nil, err
	}
	return &Nav{
		Majorcategory: majorcategory,
		Locations:     locations,
		Sublocation:   sublocaions,
	}, nil
}
