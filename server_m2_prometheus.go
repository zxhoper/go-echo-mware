package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Process is the middleware function.
func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		// TODO: something
		return nil
	}
}

// Handle is the endpoint to get stats.
func Handle(c echo.Context) error {
	s := "Hit stats endpoint"
	return c.JSON(http.StatusOK, s)
}

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")

		/////////// m1 add a custom header /////////////
		c.Response().Header().Set("MyOwnHeader", "Homerun")

		return next(c)
	}
}

func main() {
	e := echo.New()

	// Debug mode
	e.Debug = true

	//-------------------
	// Custom middleware
	//-------------------
	// Stats
	e.Use(Process)
	e.GET("/stats", Handle) // Endpoint to main process

	// Server header
	e.Use(ServerHeader)

	// Handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
