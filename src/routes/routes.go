package routes

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/estate/src/api/blog"
	"github.com/myrachanto/estate/src/api/category"
	"github.com/myrachanto/estate/src/api/dashboard"
	"github.com/myrachanto/estate/src/api/feature"
	"github.com/myrachanto/estate/src/api/gift"
	"github.com/myrachanto/estate/src/api/location"
	"github.com/myrachanto/estate/src/api/majorcategory"
	"github.com/myrachanto/estate/src/api/pages"
	"github.com/myrachanto/estate/src/api/product"
	"github.com/myrachanto/estate/src/api/profile"
	profilesections "github.com/myrachanto/estate/src/api/profileSections"
	"github.com/myrachanto/estate/src/api/seo"
	subLocation "github.com/myrachanto/estate/src/api/sublocation"
	"github.com/myrachanto/estate/src/api/tags"
	"github.com/myrachanto/estate/src/api/users"
	m "github.com/myrachanto/estate/src/middleware"

	// middle "github.com/myrachanto/estate/src/middleware"

	docs "github.com/myrachanto/estate/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ApiLoader() {
	u := users.NewUserController(users.NewUserService(users.NewUserRepo()))
	cat := category.NewcategoryController(category.NewcategoryService(category.NewcategoryRepo()))
	prof := profile.NewProfileController(profile.NewprofileService(profile.NewProfileRepo()))
	p := pages.NewpageController(pages.NewpageService(pages.NewpageRepo()))
	t := tags.NewtagController(tags.NewtagService(tags.NewtagRepo()))
	d := dashboard.NewdashboardController(dashboard.NewdashboardService(dashboard.NewdashboardRepo()))
	profSec := profilesections.NewprofileSectionsController(profilesections.NewProfilesectionsService(profilesections.NewProfileSectionRepo()))
	gi := gift.NewGiftController(gift.NewGiftService(gift.NewGiftRepo()))

	bg := blog.NewBlogController(blog.NewBlogService(blog.NewBlogRepo()))
	prod := product.NewProductController(product.NewProductService(product.NewProductRepo()))
	major := majorcategory.NewmajorcategoryController(majorcategory.NewmajorcategoryService(majorcategory.NewmajorcategoryRepo()))
	loc := location.NewlocationController(location.NewlocationService(location.NewlocationRepo()))
	subloc := subLocation.NewsubLocationController(subLocation.NewsubLocationService(subLocation.NewsubLocationRepo()))
	feat := feature.NewfeatureController(feature.NewfeatureService(feature.NewfeatureRepo()))
	seos := seo.NewseoController(seo.NewseoService(seo.NewseoRepo()))

	// bgSec := BlogSections.NewblogSectionsController(blogsections.NewblogsectionsService(blogsections.NewblogSectionRepo()))
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}
	docs.SwaggerInfo.BasePath = "/"
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = HttpErrorHandler

	api := e.Group("/api")
	// api := e.Group("/api", m.CustomAuthMidleware("admin"))
	e.Static("/", "src/public")

	{
		// auth urls
		e.POST("/register", u.Create)
		e.POST("/login", u.Login)
		e.POST("/forget", u.Forgot)

		// website urls
		e.GET("/home", d.Index)
		e.GET("/layout", d.Layout)
		e.GET("/product/category/:type", prod.GetProductsbycategory)
		e.GET("/product/properties", prod.GetProperties)
		e.GET("/product/town/:town", prod.GetOneE1)
		e.GET("/product/land/:town", prod.GetLand)
		e.GET("/product/rental/:town", prod.GetRental)
		e.GET("/product/properties/town/:sublocation", prod.GetOneE2)
		e.GET("/product/properties/:town", prod.GetProperty)
		e.GET("/product/property-type/:town", prod.GetPropertyType)
		e.GET("/search", prod.Search)
		e.GET("/search2", prod.Search2)
		e.GET("/locats", loc.GetAll)

		e.GET("/category/:category", prof.GetByCategory)
		e.GET("/navs", cat.GetAll)
		e.GET("/pages/:code", p.GetOneByUrl)
		e.POST("/comment/create", prof.CreateComment)
		e.POST("/comment/delete", prof.DeleteComment)
		e.GET("/products", prod.GetAll)
		e.GET("/products/:code", prod.GetOneE)
		e.GET("/productlimit/four", prod.GetThree)
		e.GET("/featuredproduct/four", prod.GetFeatured)
		e.GET("/products/all", prod.GetAll)
		e.GET("/product/majorcategory/:type", prod.GetProductsbyMajorcategory)
		e.GET("/product/newarrivals", prod.GetProductsbyarrival)
		e.GET("/product/hotdeals", prod.GetProductshotdeals)
		// e.GET("/product/related/:code", prod.GetOneE1)

		e.GET("/navs/all", prod.GetNavs)

		// blogs
		e.GET("/blogs/all", bg.GetAll)
		e.GET("/blogs/:url", bg.GetOneByUrl)

		// health status url
		e.GET("/health", HealthCheck).Name = "health"

		// urls documentation
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// dashboard url
	api.GET("/dashboard", d.HomeCms)

	// protected user urls
	api.GET("/logout", u.Logout, m.CustomAuthMidleware("admin"))
	api.POST("/users/shop", u.Create, m.CustomAuthMidleware("admin"))
	api.GET("/users", u.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/users/:code", u.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/users/:code", u.Update, m.CustomAuthMidleware("admin"))
	api.PUT("/users/password", u.PasswordUpdate, m.CustomAuthMidleware("admin"))
	api.PUT("/admin/update/:code", u.UpdateAdmin, m.CustomAuthMidleware("admin"))
	// api.PUT("/admin/update/:code", u.Updateadmin, m.CustomAuthMidleware("admin"))

	// category urls
	api.POST("/categorys", cat.Create, m.CustomAuthMidleware("admin"))
	api.GET("/categorys", cat.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/allcategorys", cat.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/categorys/:code", cat.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/categorys/:code", cat.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/categorys/:code", cat.Delete, m.CustomAuthMidleware("admin"))

	// features urls
	api.POST("/features", feat.Create, m.CustomAuthMidleware("admin"))
	api.GET("/features", feat.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/allfeatures", feat.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/features/:code", feat.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/features/:code", feat.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/features/:code", feat.Delete, m.CustomAuthMidleware("admin"))
	/////////////majorcategory//////////////////////////////////
	api.POST("/majorcategorys", major.Create, m.CustomAuthMidleware("admin"))
	api.GET("/majorcategorys", major.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/allmajors", major.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/majorcategorys/:code", major.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/majorcategorys/:code", major.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/majorcategorys/:code", major.Delete, m.CustomAuthMidleware("admin"))

	// location urls
	api.POST("/locations", loc.Create, m.CustomAuthMidleware("admin"))
	api.GET("/locations", loc.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/allLocations", loc.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/locations/:code", loc.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/locations/:code", loc.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/locations/:code", loc.Delete, m.CustomAuthMidleware("admin"))

	// sub locations  urls
	api.POST("/subLocations", subloc.Create, m.CustomAuthMidleware("admin"))
	api.GET("/subLocations", subloc.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/allsubLocations", subloc.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/getByLocation/:code", subloc.GetAllByLocation, m.CustomAuthMidleware("admin"))
	api.GET("/subLocations/:code", subloc.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/subLocations/:code", subloc.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/subLocations/:code", subloc.Delete, m.CustomAuthMidleware("admin"))

	// seos urls
	api.POST("/seos", seos.Create, m.CustomAuthMidleware("admin"))
	api.GET("/seos", seos.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/allseos", seos.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/seos/:code", seos.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/seos/:code", seos.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/seos/:code", seos.Delete, m.CustomAuthMidleware("admin"))
	// tags urls
	api.POST("/tags", t.Create, m.CustomAuthMidleware("admin"))
	api.GET("/tags", t.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/alltags", t.GetAll1, m.CustomAuthMidleware("admin"))
	api.GET("/tags/:code", t.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/tags/:code", t.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/tags/:code", t.Delete, m.CustomAuthMidleware("admin"))

	// pages urls
	api.POST("/pages", p.Create, m.CustomAuthMidleware("admin"))
	api.GET("/pages", p.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/pages/:code", p.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/pages/:code", p.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/pages/:code", p.Delete, m.CustomAuthMidleware("admin"))
	// pages urls
	api.POST("/gifts", gi.Create, m.CustomAuthMidleware("admin"))
	api.GET("/gifts", gi.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/gifts/:code", gi.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/gifts/:code", gi.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/gifts/:code", gi.Delete, m.CustomAuthMidleware("admin"))

	////////////////////////////////////////////////////////
	/////////////products//////////////////////////////////
	api.POST("/products", prod.Create, m.CustomAuthMidleware("admin"))
	api.GET("/products", prod.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/results/:style", prod.Results, m.CustomAuthMidleware("admin"))
	api.GET("/products/:code", prod.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/updatefeatured/:code", prod.UpdateFeatured, m.CustomAuthMidleware("admin"))
	api.PUT("/updatehotdeals/:code", prod.UpdateHotdeals, m.CustomAuthMidleware("admin"))
	api.PUT("/updatepromotion/:code", prod.UpdatePromotion, m.CustomAuthMidleware("admin"))
	api.PUT("/updateComplete/:code", prod.UpdateCompleted, m.CustomAuthMidleware("admin"))
	api.PUT("/UpdateSold/:code", prod.UpdateSold, m.CustomAuthMidleware("admin"))
	api.PUT("/likes/:code", prod.Likes, m.CustomAuthMidleware("admin"))
	api.PUT("/products/update/:code", prod.Update, m.CustomAuthMidleware("admin"))
	api.PUT("/inventory/:code", prod.AUpdate, m.CustomAuthMidleware("admin"))
	api.DELETE("/products/:code", prod.Delete, m.CustomAuthMidleware("admin"))
	// news urls
	// news urls
	api.POST("/blogs", bg.Create, m.CustomAuthMidleware("admin"))
	api.GET("/blogs", bg.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/blogs/:code", bg.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/blogs/trending/:code", bg.UpdateTrending, m.CustomAuthMidleware("admin"))
	api.PUT("/blogs/exclusive/:code", bg.UpdateExclusive, m.CustomAuthMidleware("admin"))
	api.PUT("/blogs/featured/:code", bg.UpdateFeatured, m.CustomAuthMidleware("admin"))
	api.PUT("/blogs/:code", bg.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/blogs/:code", bg.Delete, m.CustomAuthMidleware("admin"))

	//newssections
	api.PUT("/newsSections/:code", profSec.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/newsSections/:code", profSec.Delete, m.CustomAuthMidleware("admin"))
	//blogssections
	// api.PUT("/blogssections/:code", bgSec.Update, m.CustomAuthMidleware("admin"))
	// api.DELETE("/blogssections/:code", bgSec.Delete, m.CustomAuthMidleware("admin"))

	PORT := os.Getenv("PORT")
	// log.Println("fired up .... on port :1200")

	e.Logger.Fatal(e.Start(PORT))
}
