package routes

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/sports/src/api/category"
	"github.com/myrachanto/sports/src/api/dashboard"
	"github.com/myrachanto/sports/src/api/news"
	newssections "github.com/myrachanto/sports/src/api/newsSections"
	"github.com/myrachanto/sports/src/api/pages"
	"github.com/myrachanto/sports/src/api/tags"
	"github.com/myrachanto/sports/src/api/users"
	m "github.com/myrachanto/sports/src/middleware"

	docs "github.com/myrachanto/sports/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ApiLoader() {
	u := users.NewUserController(users.NewUserService(users.NewUserRepo()))
	cat := category.NewcategoryController(category.NewcategoryService(category.NewcategoryRepo()))
	n := news.NewnewsController(news.NewnewsService(news.NewnewsRepo()))
	p := pages.NewpageController(pages.NewpageService(pages.NewpageRepo()))
	t := tags.NewtagController(tags.NewtagService(tags.NewtagRepo()))
	d := dashboard.NewdashboardController(dashboard.NewdashboardService(dashboard.NewdashboardRepo()))
	nesecs := newssections.NewnewsSectionsController(newssections.NewnewssectionsService(newssections.NewnewsSectionRepo()))
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

	api := e.Group("/api", m.CustomAuthMidleware("users"))
	e.Static("/", "src/public")

	{
		// auth urls
		e.POST("/register", u.Create)
		e.POST("/login", u.Login)
		e.POST("/forget", u.Forgot)

		// website urls
		e.GET("/home", d.Index)
		e.GET("/news", n.GetAll)
		e.GET("/category/:category", n.GetByCategory)
		e.GET("/news/:url", n.GetOneByUrl)
		e.GET("/navs", cat.GetAll)
		e.GET("/pages/:code", p.GetOneByUrl)

		// health status url
		e.GET("/health", HealthCheck).Name = "health"

		// urls documentation
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// dashboard url
	api.GET("/dashboard", d.HomeCms)

	// protected user urls
	api.GET("/logout", u.Logout)
	api.POST("/users/shop", u.Create)
	api.GET("/users", u.GetAll)
	api.GET("/users/:code", u.GetOne)
	api.PUT("/users/:code", u.Update)
	api.PUT("/users/password", u.PasswordUpdate)
	api.PUT("/admin/update/:code", u.UpdateAdmin)

	// category urls
	api.POST("/category", cat.Create)
	api.GET("/category", cat.GetAll)
	api.GET("/category/:code", cat.GetOne)
	api.PUT("/category/:code", cat.Update)
	api.DELETE("/category/:code", cat.Delete)

	// tags urls
	api.POST("/tags", t.Create)
	api.GET("/tags", t.GetAll)
	api.GET("/tags/:code", t.GetOne)
	api.PUT("/tags/:code", t.Update)
	api.DELETE("/tags/:code", t.Delete)

	// pages urls
	api.POST("/pages", p.Create)
	api.GET("/pages", p.GetAll)
	api.GET("/pages/:code", p.GetOne)
	api.PUT("/pages/:code", p.Update)
	api.DELETE("/pages/:code", p.Delete)

	// news urls
	api.POST("/news", n.Create)
	api.GET("/news", n.GetAll)
	api.GET("/news/:code", n.GetOne)
	api.PUT("/news/trending/:code", n.UpdateTrending)
	api.PUT("/news/exclusive/:code", n.UpdateExclusive)
	api.PUT("/news/featured/:code", n.UpdateFeatured)
	api.PUT("/news/:code", n.Update)
	api.DELETE("/news/:code", n.Delete)

	//newssections
	api.PUT("/newsSections/:code", nesecs.Update)
	api.DELETE("/newsSections/:code", nesecs.Delete)

	PORT := os.Getenv("PORT")
	// log.Println("fired up .... on port :1200")

	e.Logger.Fatal(e.Start(PORT))
}
