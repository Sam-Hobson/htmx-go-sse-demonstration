package main

import (
	"salad2/internal/routes"
	"salad2/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = utils.Template

	// Expose public files from server
	e.Static("/public", "public")

	// Setup routes
	e.GET("/", routes.GetRootRoute())
	e.GET("/contacts", routes.GetContactsRoute())
	e.POST("/contacts", routes.PostContactsRoute())
	e.POST("/delete", routes.PostDeleteReoute())

	e.Logger.Fatal(e.Start(":8080"))
}
