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
	api.GET("/logout", u.Logout, m.CustomAuthMidleware("admin"))
	api.POST("/users/shop", u.Create, m.CustomAuthMidleware("admin"))
	api.GET("/users", u.GetAll, m.CustomAuthMidlewareAuditor("admin"))
	api.GET("/users/:code", u.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/users/:code", u.Update, m.CustomAuthMidleware("admin"))
	api.PUT("/users/password", u.PasswordUpdate, m.CustomAuthMidleware("admin"))
	api.PUT("/admin/update/:code", u.UpdateAdmin, m.CustomAuthMidleware("admin"))
	api.PUT("/auditor/update/:code", u.UpdateAuditor, m.CustomAuthMidleware("admin"))

	// category urls
	api.POST("/category", cat.Create, m.CustomAuthMidleware("admin"))
	api.GET("/category", cat.GetAll, m.CustomAuthMidlewareAuditor("admin"))
	api.GET("/category/:code", cat.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/category/:code", cat.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/category/:code", cat.Delete, m.CustomAuthMidleware("admin"))

	// tags urls
	api.POST("/tags", t.Create, m.CustomAuthMidleware("admin"))
	api.GET("/tags", t.GetAll, m.CustomAuthMidlewareAuditor("admin"))
	api.GET("/tags/:code", t.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/tags/:code", t.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/tags/:code", t.Delete, m.CustomAuthMidleware("admin"))

	// pages urls
	api.POST("/pages", p.Create, m.CustomAuthMidleware("admin"))
	api.GET("/pages", p.GetAll, m.CustomAuthMidleware("admin"))
	api.GET("/pages/:code", p.GetOne, m.CustomAuthMidleware("admin"))
	api.PUT("/pages/:code", p.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/pages/:code", p.Delete, m.CustomAuthMidleware("admin"))

	// news urls
	api.POST("/news", n.Create, m.CustomAuthMidlewareAuditor("auditor"))
	api.GET("/news", n.GetAll, m.CustomAuthMidlewareAuditor("auditor"))
	api.GET("/news/:code", n.GetOne, m.CustomAuthMidlewareAuditor("auditor"))
	api.PUT("/news/trending/:code", n.UpdateTrending, m.CustomAuthMidlewareAuditor("admin"))
	api.PUT("/news/exclusive/:code", n.UpdateExclusive, m.CustomAuthMidlewareAuditor("admin"))
	api.PUT("/news/featured/:code", n.UpdateFeatured, m.CustomAuthMidlewareAuditor("admin"))
	api.PUT("/news/:code", n.Update, m.CustomAuthMidlewareAuditor("auditor"))
	api.DELETE("/news/:code", n.Delete, m.CustomAuthMidleware("admin"))

	//newssections
	api.PUT("/newsSections/:code", nesecs.Update, m.CustomAuthMidleware("admin"))
	api.DELETE("/newsSections/:code", nesecs.Delete, m.CustomAuthMidleware("admin"))

	PORT := os.Getenv("PORT")
	// log.Println("fired up .... on port :1200")

	e.Logger.Fatal(e.Start(PORT))
}
