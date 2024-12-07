package handlers

import (
	"log"
	"net/http"
	"shopApi/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetOrders(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID из параметров маршрута
		id := c.Param("user_id")
		log.Println("Полученный параметр idStr:", id)

		// Убираем лишние пробелы
		id = strings.TrimSpace(id)

		// Выполняем запрос к базе данных
		var orders []models.Order
		err := db.Select(&orders, "SELECT * FROM orders WHERE user_id = $1", id)

		if err != nil {
			log.Println("Ошибка запроса к базе данных:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заказов"})
			return
		}

		// Отправляем ответ
		c.JSON(http.StatusOK, orders)
	}
}

func GetOrderItems(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID из параметров маршрута
		userId := c.Param("user_id")
		orderIdStr := c.Param("order_id")

		// Убираем лишние пробелы
		orderIdStr = strings.TrimSpace(orderIdStr)
		userId = strings.TrimSpace(userId)

		// Преобразуем ID в целое число
		order_id, err := strconv.Atoi(orderIdStr)
		if err != nil {
			log.Println("Ошибка преобразования idStr в int:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
			return
		}

		// Выполняем запрос к базе данных
		var orderItems []models.Product
		err = db.Select(&orderItems, `
			SELECT 
				pr.product_id,
				"name",
				description,
				price,
				stock,
				image_url,
				EXISTS (
						SELECT 1 
						FROM Favorites f 
						WHERE f.product_id = pr.product_id AND f.user_id = $1
					) AS is_favorite
				FROM order_items oi
				JOIN product pr ON pr.product_id = oi.product_id
				WHERE order_id = $2;`,
			userId, order_id)

		if err != nil {
			log.Println("Ошибка запроса к базе данных:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заказов"})
			return
		}

		// Отправляем ответ
		c.JSON(http.StatusOK, orderItems)
	}
}

func CreateOrder(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		log.Println("Received POST request on /orders/:user_id")

		// Привязываем данные из тела запроса
		if err := c.ShouldBindJSON(&order); err != nil {
			log.Println("ShouldBindJSON failed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка инициализации транзакции"})
			return
		}

		// Вставляем заказ в таблицу orders
		queryOrder := `
			INSERT INTO orders (user_id, total, status)
			VALUES (:user_id, :total, :status)
			RETURNING order_id
		`
		rows, err := tx.NamedQuery(queryOrder, &order)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления заказа"})
			return
		}
		if rows.Next() {
			rows.Scan(&order.OrderID) // Получаем ID нового заказа
		}
		rows.Close()

		// Вставляем товары в таблицу order_products
		queryProducts := `
			INSERT INTO order_items (order_id, product_id, quantity)
			VALUES (:order_id, :product_id, :quantity)
		`
		for _, item := range order.CartItems {
			productData := map[string]interface{}{
				"order_id":   order.OrderID,
				"product_id": item.ProductID,
				"quantity":   item.Quantity, // Здесь quantity (например, 1, 2 и т.д.) нужно передать из тела запроса
			}
			_, err := tx.NamedExec(queryProducts, productData)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления товаров к заказу"})
				return
			}
		}

		userId := c.Param("user_id")
		_, err = db.Exec("DELETE FROM cart WHERE user_id = $1", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки корзины"})
			return
		}

		// Завершаем транзакцию
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения заказа"})
			return
		}

		// Отправляем ответ
		c.JSON(http.StatusCreated, order)
	}
}
