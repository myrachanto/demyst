package middle

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/sports/src/pasetos"
)

const (
	authorisationHeaderKey = "Authorization"
	authorisationType      = "Bearer"
)

func CustomAuthMidleware(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get(authorisationHeaderKey)
			if len(authorizationHeader) == 0 {
				return c.JSON(http.StatusUnauthorized, "Authorization header not provided")
			}
			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				return c.JSON(http.StatusUnauthorized, "You Have No Authorization")
			}
			authtype := fields[0]
			if authtype != authorisationType {

				return c.JSON(http.StatusUnauthorized, "That type of Authorization is not allowed here!")
			}
			accessToken := fields[1]
			maker, _ := pasetos.NewPasetoMaker()
			payload, err := maker.VerifyToken(accessToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "That token is invalid!")
			}

			if !payload.Admin {
				return c.JSON(http.StatusUnauthorized, "You Are not Authorised!")
			}
			return next(c)
		}
	}
}
func CustomAuthMidlewareAuditor(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get(authorisationHeaderKey)
			if len(authorizationHeader) == 0 {
				return c.JSON(http.StatusUnauthorized, "Authorization header not provided")
			}
			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				return c.JSON(http.StatusUnauthorized, "You Have No Authorization")
			}
			authtype := fields[0]
			if authtype != authorisationType {

				return c.JSON(http.StatusUnauthorized, "That type of Authorization is not allowed here!")
			}
			accessToken := fields[1]
			maker, _ := pasetos.NewPasetoMaker()
			payload, err := maker.VerifyToken(accessToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "That token is invalid!")
			}
			// fmt.Println("----------------payload", payload)
			if !payload.Auditor {
				return c.JSON(http.StatusUnauthorized, "You Are not Authorised!")
			}
			return next(c)
		}
	}
}

func PasetoAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// fmt.Println("-----------------------step 1")
		authorizationHeader := c.Request().Header.Get(authorisationHeaderKey)
		if len(authorizationHeader) == 0 {
			return c.JSON(http.StatusUnauthorized, "Authorization header not provided")
		}
		// fmt.Println("-----------------------step 2")
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return c.JSON(http.StatusUnauthorized, "You Have No Authorization")
		}
		authtype := fields[0]
		if authtype != authorisationType {

			return c.JSON(http.StatusUnauthorized, "That type of Authorization is not allowed here!")
		}
		// fmt.Println("-----------------------step 3")
		accessToken := fields[1]
		maker, _ := pasetos.NewPasetoMaker()
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "That token is invalid!")
		}
		// fmt.Println("-----------------------step 4", payload.Admin)
		if !payload.Admin {
			return c.JSON(http.StatusUnauthorized, "You Are not Authorised!")
		}
		return next(c)
	}
}
