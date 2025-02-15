package routes

import (
	"net/http"
	"salad2/internal/data"
	"salad2/internal/events"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func PostDeleteReoute() echo.HandlerFunc {
	return func(c echo.Context) error {
		mode := strings.ToUpper(c.QueryParam("mode"))
		if mode != "" && mode != "JSON" && mode != "HTML" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid value for query parameter \"mode\"")
		}

		idStr := c.QueryParam("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			if mode == "HTML" {
				return c.HTML(http.StatusBadRequest, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		contactToDelete, err := data.MainData.ContactList.FindFunc(func(c *events.Contact) bool { return c.Id == id })
		if err != nil {
			if mode == "HTML" {
				return c.HTML(http.StatusBadRequest, "Invalid id")
			}
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}

		events.RemoveContactQueue.Queue(*contactToDelete)

		if mode == "HTML" {
			return c.HTML(http.StatusOK, "")
		}
		return c.JSON(http.StatusOK, contactToDelete)
	}
}
