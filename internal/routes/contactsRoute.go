package routes

import (
	"context"
	"fmt"
	"net/http"
	"salad2/internal/data"
	"salad2/internal/events"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func PostContactsRoute() echo.HandlerFunc {
	return func(c echo.Context) error {
		mode := strings.ToUpper(c.QueryParam("mode"))
		if mode != "JSON" && mode != "HTML" && mode != "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid value for query parameter \"mode\"")
		}

		contact := events.NewContact()

		if err := c.Bind(contact); err != nil {
			if mode == "HTML" {
				return c.HTML(http.StatusBadRequest, "Invalid data format")
			}
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data format")
		}
		if err := contact.Validate(); err != nil {
			if mode == "HTML" {
				return c.HTML(http.StatusBadRequest, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		events.ContactQueue.Queue(*contact)

		if mode == "HTML" {
			return c.HTML(http.StatusOK, "")
		}
		return c.JSON(http.StatusOK, contact)
	}
}

func GetContactsRoute() echo.HandlerFunc {
	return func(c echo.Context) error {
		mode := strings.ToUpper(c.QueryParam("mode"))
		if mode != "SSE" && mode != "JSON" && mode != "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid value for query parameter \"mode\"")
		}

		if mode == "SSE" {
			uniqueTopic := fmt.Sprintf("GetContactsRoute-%s", uuid.New().String())

			addContactEvents := events.ContactQueue.Subscribe(uniqueTopic, 100)
			removeContactEvents := events.RemoveContactQueue.Subscribe(uniqueTopic, 100)
			defer events.ContactQueue.Unsubscribe(uniqueTopic)
			defer events.RemoveContactQueue.Unsubscribe(uniqueTopic)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			addContactProvider := EventDataProviderFromTemplate(ctx, addContactEvents, "add-contact", "contactRow")
			removeContactProvider := EventDataProviderFromTemplate(ctx, removeContactEvents, "delete-contact", "deleteContactRow")

			return HandleSSEConnection(c, addContactProvider, removeContactProvider)
		}

		return c.JSON(http.StatusOK, data.MainData.ContactList.Items())
	}
}
