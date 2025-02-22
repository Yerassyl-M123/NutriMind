package main

import (
	"project/config"
	"project/internal/handlers"
	"project/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	server := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 5 * 60})
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.SessionTimeout())

	server.LoadHTMLGlob("templates/*")

	server.GET("/signUpPage", handlers.SignUpPage)
	server.POST("/signup", handlers.SignUp)
	server.GET("/signInPage", handlers.SignInPage)
	server.POST("/signin", handlers.SignIn)
	server.GET("/logout", handlers.Logout)

	server.Run(":8080")
}
