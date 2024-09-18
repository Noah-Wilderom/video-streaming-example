package handler

import (
	"context"
	"fmt"
	"github.com/Noah-Wilderom/video-streaming-test/client"
	"github.com/Noah-Wilderom/video-streaming-test/pkg/kit/validate"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views/auth"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func HandleLoginIndex(c echo.Context) error {
	return Render(c, auth.Login())
}

func HandleLoginCreate(c echo.Context) error {
	credentials := auth.UserCredentials{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	errors := auth.LoginErrors{}
	if ok := validate.New(&credentials, validate.Fields{
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Required, validate.Min(8)),
	}).Validate(&errors); !ok {
		return Render(c, auth.LoginForm(credentials, errors))
	}

	apiClient := client.NewClient()

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	resp, code, err := apiClient.Login(ctx, credentials.Email, credentials.Password)
	if err != nil {
		fmt.Println("api", err.Error())
		if code == http.StatusUnauthorized || code == http.StatusBadRequest {
			errors.InvalidCredentials = err.Error()
			return Render(c, auth.LoginForm(credentials, errors))
		}

		return err
	}

	setAuthCookie(c, resp.Token)
	cookie, err := c.Cookie("intended")
	if err == nil && len(cookie.Value) > 0 {
		intended := cookie.Value
		cookie.MaxAge = -1
		cookie.Value = ""
		c.SetCookie(cookie)
		return HxRedirect(c, intended)
	}

	return HxRedirect(c, "/")
}

func HandleSignupIndex(c echo.Context) error {
	return Render(c, auth.Signup())
}

func HandleSignupCreate(c echo.Context) error {
	params := auth.SignupParams{
		Name:            c.FormValue("name"),
		Email:           c.FormValue("email"),
		Password:        c.FormValue("password"),
		ConfirmPassword: c.FormValue("confirm_password"),
	}

	errors := auth.SignupErrors{}
	if ok := validate.New(&params, validate.Fields{
		"Name":     validate.Rules(validate.Required),
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Required, validate.Min(8)),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.Password),
			validate.Message("passwords do not match"),
		),
	}).Validate(&errors); !ok {
		return Render(c, auth.SignupForm(params, errors))
	}

	apiClient := client.NewClient()

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	resp, code, err := apiClient.Signup(ctx, params.Name, params.Email, params.Password)

	if err != nil {
		if code == http.StatusUnauthorized || code == http.StatusBadRequest {
			fmt.Println("Error code:", code)
		}

		return err
	}

	setAuthCookie(c, resp.Token)
	cookie, err := c.Cookie("intended")
	if err == nil && len(cookie.Value) > 0 {
		intended := cookie.Value
		cookie.MaxAge = -1
		cookie.Value = ""
		c.SetCookie(cookie)
		return HxRedirect(c, intended)
	}

	return HxRedirect(c, "/")
}

func HandleLogoutCreate(c echo.Context) error {
	cookie := http.Cookie{
		Value:    "",
		Name:     "at",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}

	c.SetCookie(&cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}

func HandleAuthCallback(c echo.Context) error {
	accessToken := c.Request().URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return Render(c, auth.CallbackScript())
	}

	setAuthCookie(c, accessToken)
	return c.Redirect(http.StatusSeeOther, "/")
}

func setAuthCookie(c echo.Context, accessToken string) {
	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	c.SetCookie(cookie)
}
