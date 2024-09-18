package main

import (
	"github.com/Noah-Wilderom/video-streaming-test/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(handler.WithGlobalData)

	e.StaticFS("/public/*", os.DirFS("public"))

	guest := e.Group("")
	guest.GET("/login", handler.HandleLoginIndex)
	guest.POST("/login", handler.HandleLoginCreate)
	guest.GET("/signup", handler.HandleSignupIndex)
	guest.POST("/signup", handler.HandleSignupCreate)
	guest.GET("/auth/callback", handler.HandleAuthCallback)

	authenticated := e.Group("", handler.Authenticated)
	authenticated.POST("/auth/logout", handler.HandleLogoutCreate)

	authenticated.GET("/", handler.HandleHomeIndex)

	authenticated.GET("/video/:id", handler.HandleVideoShow)
	authenticated.POST("/upload", handler.HandleUploadVideo)

	e.Logger.Fatal(e.Start(":3000"))
}
