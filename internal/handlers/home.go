package handlers

import (
	"net/http"
	"project/config"
	"project/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddToMealPlan(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	recipeID := c.PostForm("recipe_id")
	mealType := c.PostForm("meal_type")
	mealDate := c.PostForm("meal_date")

	_, err := config.DB.Exec("INSERT INTO meal_plan (user_id, recipe_id, meal_type, meal_date) VALUES ($1, $2, $3, $4)",
		userID, recipeID, mealType, mealDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в план"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Блюдо добавлено в план!"})
}

func GetMealPlan(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	mealDate := c.Query("date")

	rows, err := config.DB.Query(`
        SELECT recipes.id, recipes.name, meal_plan.meal_type 
        FROM meal_plan 
        JOIN recipes ON meal_plan.recipe_id = recipes.id 
        WHERE meal_plan.user_id = $1 AND meal_plan.meal_date = $2
    `, userID, mealDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения плана"})
		return
	}
	defer rows.Close()

	mealPlan := map[string][]models.Recipe{
		"breakfast": {},
		"lunch":     {},
		"dinner":    {},
		"snacks":    {},
	}

	for rows.Next() {
		var r models.Recipe
		var mealType string
		if err := rows.Scan(&r.ID, &r.Name, &mealType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
			return
		}
		mealPlan[mealType] = append(mealPlan[mealType], r)
	}

	c.JSON(http.StatusOK, mealPlan)
}

func GetHomePage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/signInPage")
		return
	}

	c.HTML(http.StatusOK, "home.html", nil)
}
