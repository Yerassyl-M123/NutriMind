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
	server.Static("/uploads", "./uploads")

	server.GET("/recipes", middleware.AuthRequired(handlers.GetRecipes))
	server.GET("/recipe/:id", middleware.AuthRequired(handlers.GetRecipeDetails))

	server.GET("/", middleware.AuthRequired(handlers.GetHomePage))
	server.POST("/add-to-plan", middleware.AuthRequired(handlers.AddToMealPlan))
	server.GET("/get-meal-plan", middleware.AuthRequired(handlers.GetMealPlan))

	server.GET("/my-products", middleware.AuthRequired(handlers.GetMyProducts))
	server.POST("/upload-product", middleware.AuthRequired(handlers.UploadProduct))
	server.POST("/search-recipes", middleware.AuthRequired(handlers.SearchRecipes))

	server.GET("/signUpPage", handlers.SignUpPage)
	server.POST("/signup", handlers.SignUp)
	server.GET("/signInPage", handlers.SignInPage)
	server.POST("/signin", handlers.SignIn)
	server.GET("/logout", handlers.Logout)

	server.Run(":8080")
}
