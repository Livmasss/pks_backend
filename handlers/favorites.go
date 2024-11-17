package handlers

import (
	"net/http"
	"shopApi/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetFavorites(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var products []models.Product

		err := db.Select(
			&products,
			`SELECT p.product_id, name, description, price, stock, image_url FROM favorites f
			JOIN product p ON f.product_id = p.product_id
			WHERE user_id = $1`,
			userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка получения списка продуктов",
				"details": err.Error(), // Логирование деталей ошибки
			})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}
func AddToFavorites(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var item struct {
			ProductID int `json:"product_id"`
		}
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}
		_, err := db.Exec("INSERT INTO Favorites (user_id, product_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			userId, item.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в избранное"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в избранное"})
	}
}

func RemoveFromFavorites(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		productId := c.Param("productId")
		_, err := db.Exec("DELETE FROM Favorites WHERE user_id = $1 AND product_id = $2", userId, productId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления из избранного"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Товар удален из избранного"})
	}
}
