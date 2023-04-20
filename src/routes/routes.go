package routes

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/demyst/src/api/users"

	docs "github.com/myrachanto/demyst/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ApiLoader() {
	u := users.NewUserController(users.NewUserService(users.NewUserRepo()))
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

	{
		e.POST("/register", u.Create)
		e.POST("/login", u.Login)
		e.GET("/health", HealthCheck).Name = "health"
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	api.GET("/logout", u.Logout)
	api.POST("/users/shop", u.Create)
	api.GET("/users", u.GetAll)
	api.GET("/users/:code", u.GetOne)
	api.PUT("/users/password", u.PasswordUpdate)

	PORT := os.Getenv("PORT")
	// log.Println("fired up .... on port :1200")

	e.Logger.Fatal(e.Start(PORT))
}
