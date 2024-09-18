package models

const UserContextKey = "user"

type AuthenticatedUser struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Token      string `json:"token,omitempty"`
	IsLoggedIn bool   `json:"is_logged_in,omitempty"`
}
