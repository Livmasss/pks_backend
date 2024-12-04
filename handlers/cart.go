package handlers

import (
	"net/http"
	"shopApi/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var cartItems []models.CartItem
		err := db.Select(&cartItems, `SELECT 
				p.product_id, p.name, p.description, p.price, p.stock, p.image_url, c.quantity
			FROM
				Cart c
			JOIN
				Product p ON c.product_id = p.product_id
			WHERE
				c.user_id = $1;`, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка получения списка продуктов",
				"details": err.Error(), // Логирование деталей ошибки
			})
			return
		}

		c.JSON(http.StatusOK, cartItems)
	}
}

func AddToCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var item struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}
		_, err := db.Exec("INSERT INTO Cart (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = Cart.quantity + $3",
			userId, item.ProductID, item.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в корзину"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в корзину"})
	}
}

func DecreaseCartQuantity(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var item struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}
		_, err := db.Exec(
			`UPDATE Cart
			SET quantity = Cart.quantity - $3
			WHERE user_id = $1 AND product_id = $2;`,
			userId, item.ProductID, item.Quantity)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка уменьшения количества"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Количество уменьшено"})
	}
}

func RemoveFromCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		productId := c.Param("productId")
		_, err := db.Exec("DELETE FROM Cart WHERE user_id = $1 AND product_id = $2", userId, productId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления из корзины"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Товар удален из корзины"})
	}
}

func ClearCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		_, err := db.Exec("DELETE FROM Cart WHERE user_id = $1", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки корзины"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Корзина очищена"})
	}
}
