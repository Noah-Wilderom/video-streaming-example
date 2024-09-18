package handler

import (
	"github.com/Noah-Wilderom/video-streaming-test/client"
	"github.com/Noah-Wilderom/video-streaming-test/models"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

func WithGlobalData(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.Contains(c.Request().URL.Path, "/public") {
			return next(c)
		}

		data := views.GlobalData{
			Name:        os.Getenv("APP_NAME"),
			Host:        os.Getenv("APP_HOST"),
			Version:     os.Getenv("APP_VERSION"),
			Environment: os.Getenv("APP_ENV"),
			Debug:       os.Getenv("APP_DEBUG") == "true",
		}

		c.Set(views.GlobalDataContextKey, data)

		return next(c)
	}
}

func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.Contains(c.Request().URL.Path, "/public") {
			return next(c)
		}

		cookie, err := c.Cookie("at")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
			return next(c)
		}

		apiClient := client.NewClient()
		checkResponse, err := apiClient.Check(c.Request().Context(), cookie.Value)
		if err != nil {
			invalidateCookie := http.Cookie{
				Value:    "",
				Name:     "at",
				MaxAge:   -1,
				HttpOnly: true,
				Path:     "/",
				Secure:   true,
			}

			c.SetCookie(&invalidateCookie)
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		user := models.AuthenticatedUser{
			Name:       checkResponse.User.Name,
			Email:      checkResponse.User.Email,
			Token:      checkResponse.Token,
			IsLoggedIn: true,
		}

		c.Set(models.UserContextKey, user)

		return next(c)
	}
}
