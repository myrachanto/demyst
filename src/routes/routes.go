package routes

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/demyst/src/api/accounting"
	"github.com/myrachanto/demyst/src/api/loan"
	"github.com/myrachanto/demyst/src/api/users"
	m "github.com/myrachanto/demyst/src/middleware"

	docs "github.com/myrachanto/demyst/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ApiLoader() {
	u := users.NewUserController(users.NewUserService(users.NewUserRepo()))
	l := loan.NewloanController(loan.NewloanService(loan.NewloanRepo()))
	a := accounting.NewaccountingController(accounting.NewaccountingService(accounting.NewaccountingRepo()))
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
	e.Static("/", "src/public")

	{
		e.POST("/register", u.Create)
		e.POST("/login", u.Login)
		e.GET("/health", HealthCheck).Name = "health"
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	api.GET("/logout", u.Logout)
	api.POST("/users/shop", u.Create)
	api.GET("/users", u.GetAll, m.PasetoAuthMiddleware)
	api.GET("/users/:code", u.GetOne, m.PasetoAuthMiddleware)
	api.PUT("/users/password", u.PasswordUpdate, m.PasetoAuthMiddleware)

	api.POST("/loans", l.Create, m.PasetoAuthMiddleware)
	api.GET("/loans", l.GetAll, m.PasetoAuthMiddleware)
	api.GET("/loans/:code", l.GetOne, m.PasetoAuthMiddleware)
	api.PUT("/loans/:code", l.Submit, m.PasetoAuthMiddleware)

	api.POST("/accountings", a.Create, m.PasetoAuthMiddleware)
	api.GET("/accountings", a.GetAll, m.PasetoAuthMiddleware)
	api.GET("/accountings/:code", a.GetOne, m.PasetoAuthMiddleware)
	api.PUT("/accountings/:code", a.Update, m.PasetoAuthMiddleware)
	api.DELETE("/accountings/:code", a.Delete, m.PasetoAuthMiddleware)

	PORT := os.Getenv("PORT")
	// log.Println("fired up .... on port :1200")

	e.Logger.Fatal(e.Start(PORT))
}
