package handlers

import (
	"encoding/json"
	"net/http"
	"project/config"
	"project/internal/models"

	"github.com/gin-gonic/gin"
)

func GetRecipes(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, image_url FROM recipes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки рецептов"})
		return
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var r models.Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.ImageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
			return
		}
		recipes = append(recipes, r)
	}

	c.HTML(http.StatusOK, "recipes.html", gin.H{"recipes": recipes})
}

func GetRecipeDetails(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	err := config.DB.QueryRow("SELECT id, name, image_url, cooking_time, ingredients, steps FROM recipes WHERE id = $1", id).
		Scan(&recipe.ID, &recipe.Name, &recipe.ImageURL, &recipe.CookingTime, &recipe.Ingredients, &recipe.Steps)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Рецепт не найден"})
		return
	}

	var ingredients []map[string]string
	var steps []string

	if err := json.Unmarshal(recipe.Ingredients, &ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ингредиентов"})
		return
	}

	if err := json.Unmarshal(recipe.Steps, &steps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки шагов приготовления"})
		return
	}

	c.HTML(http.StatusOK, "recipe_details.html", gin.H{
		"recipe":      recipe,
		"ingredients": ingredients,
		"steps":       steps,
	})
}
