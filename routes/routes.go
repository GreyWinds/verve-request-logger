package routes

import (
	"Users/vaibhav.sabharwal/verve/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/api/verve/accept", handlers.AcceptHandler)
}
