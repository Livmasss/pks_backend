package handlers

import (
	"net/http"
	"shopApi/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}

		_, err := db.Query(`INSERT INTO users(user_id, username, email) VALUES($1, $2, $3)`, user.UserID, user.Username, user.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка получения списка продуктов",
				"details": err.Error(), // Логирование деталей ошибки
			})
			return
		}

		c.JSON(http.StatusOK, "")
	}
}

func GetProfile(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var user models.User
		err := db.Get(&user, `SELECT * FROM users WHERE user_id = $1`, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка получения профиля",
				"details": err.Error(), // Логирование деталей ошибки
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
