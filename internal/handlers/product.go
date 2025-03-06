package handlers

import (
	"fmt"
	"net/http"
	"project/config"
	"project/internal/models"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetMyProducts(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/signInPage")
		return
	}

	rows, err := config.DB.Query("SELECT id, name, photo, number FROM products WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки продуктов"})
		return
	}
	defer rows.Close()

	var products []struct {
		ID     int
		Name   string
		Photo  string
		Number string
	}

	for rows.Next() {
		var p struct {
			ID     int
			Name   string
			Photo  string
			Number string
		}
		if err := rows.Scan(&p.ID, &p.Name, &p.Photo, &p.Number); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
			return
		}
		products = append(products, p)
	}

	c.HTML(http.StatusOK, "my_products.html", gin.H{"products": products})
}

func UploadProduct(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	name := c.PostForm("name")
	number := c.PostForm("number")

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка загрузки файла"})
		return
	}

	photoPath := "uploads/" + file.Filename

	err = c.SaveUploadedFile(file, photoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения фото"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO products (user_id, name, photo, number) VALUES ($1, $2, $3, $4)",
		userID, name, photoPath, number)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения продукта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Продукт добавлен!"})
}

func SearchRecipes(c *gin.Context) {
	var request struct {
		Products []string `json:"products"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка обработки запроса"})
		return
	}

	if len(request.Products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не выбраны продукты"})
		return
	}

	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	query := `
        SELECT id, name, image_url 
        FROM recipes 
        WHERE `

	conditions := []string{}
	args := []interface{}{}

	for i, product := range request.Products {
		conditions = append(conditions, fmt.Sprintf("LOWER(ingredients::TEXT) ILIKE LOWER($%d)", i+1))
		args = append(args, "%"+product+"%")
	}

	query += strings.Join(conditions, " AND ")

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска рецептов"})
		return
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var r models.Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.ImageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
			return
		}
		recipes = append(recipes, r)
	}

	c.JSON(http.StatusOK, recipes)
}
