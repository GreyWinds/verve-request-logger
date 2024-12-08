package main

import (
	"Users/vaibhav.sabharwal/verve/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	//Register routes
	routes.RegisterRoutes(e)

	//Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
