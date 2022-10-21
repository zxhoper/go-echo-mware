package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "an_counter_myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

// MiddlewareFirst is the middleware function.
//   - Store Value to Context:
//     You can use echo.Context.Set and echo.Context.Get
//     to store value in middleware into context
//     and get that value from context in downstream
//     middlewares or handler
func MiddlewareFirst(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: something
		fmt.Printf(" Mark -> Time:%s, Hit MiddlewareFirst middleware\n", time.Now())

		///////// try echo.Context.Set and echo.Context.Get //////////
		i := 1
		// Use echo.Context.Set and echo.Context.Get
		// to store value in middleware into context
		c.Set("CounterInContext", i)
		fmt.Printf(" CounterInContext -> %d in MiddlewareFirst middleware, -> Time:%s\n",
			c.Get("CounterInContext").(int),
			time.Now())
		//\\\\\\\ try echo.Context.Set and echo.Context.Get //////////

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// MiddlewareSecond is the middleware function.
func MiddlewareSecond(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// TODO: something
		fmt.Printf(" Mark -> Time:%s, Hit MiddlewareSecond middleware\n", time.Now())

		///////// try echo.Context.Set and echo.Context.Get //////////
		i := c.Get("CounterInContext").(int)
		fmt.Printf(" CounterInContext -> %d in MiddlewareSecond middleware before increment-> Time:%s\n",
			i,
			time.Now())
		i++
		c.Set("CounterInContext", i)
		fmt.Printf(" CounterInContext -> %d in MiddlewareSecond middleware-> Time:%s\n",
			c.Get("CounterInContext").(int),
			time.Now())
		//\\\\\\\ try echo.Context.Set and echo.Context.Get //////////

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// MiddlewareAddServerHeader middleware adds a `Server` header to the response.
func MiddlewareAddServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/4.0")

		/////////// m1 add a custom header /////////////
		c.Response().Header().Set("MyOwnHeader", "Home run")

		fmt.Printf(" Mark -> Time:%s, Hit MiddlewareAddServerHeader middleware\n", time.Now())

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// myRootHandler is the main endpoint
func myRootHandler(c echo.Context) error {

	// prometheus counter INCrease
	opsProcessed.Inc()

	fmt.Printf(" Mark -> Time:%s, Hit Endpoint [myRootHandler] \n\n\n\n", time.Now())

	return c.String(http.StatusOK, "Hello, World!\n\n\n")
}

// myStatsHandle is the second endpoint
func myStatsHandle(c echo.Context) error {
	s := "Hit StatsHandle endpoint\n\n\n"
	fmt.Printf(" Mark -> Time:%s, Hit Endpoint [myStatsHandle] \n\n\n\n", time.Now())
	return c.JSON(http.StatusOK, s)
}

func main() {
	e := echo.New()

	// Debug mode
	e.Debug = true

	// Add AddServerHeader middleware to all endpoints
	e.Use(MiddlewareAddServerHeader)

	// Add MiddlewareFirst middleware to all endpoints
	e.Use(MiddlewareFirst)

	e.GET("/stats", myStatsHandle) // Second Endpoint

	// Handler
	e.GET("/", myRootHandler, MiddlewareSecond)

	// prometheus endpoint
	// -> error e.GET("/metrics", promhttp.Handler())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
