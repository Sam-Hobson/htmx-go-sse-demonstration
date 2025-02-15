package routes

import (
	"net/http"
	"salad2/internal/data"

	"github.com/labstack/echo/v4"
)

func GetRootRoute() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "root", data.MainData)
	}
}
