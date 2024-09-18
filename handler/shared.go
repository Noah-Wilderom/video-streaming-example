package handler

import (
	"github.com/Noah-Wilderom/video-streaming-test/models"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func HxRedirect(c echo.Context, to string) error {
	if c.Request().Header.Get("HX-Request") != "" {
		c.Response().Header().Set("HX-Redirect", to)
		c.Response().WriteHeader(http.StatusSeeOther)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, to)
}

func getAuthenticatedUser(c echo.Context) models.AuthenticatedUser {
	user, ok := c.Get(models.UserContextKey).(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}
	}

	return user
}
