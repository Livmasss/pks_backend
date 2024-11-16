package main

import (
	"log"
	"shopApi/db"
	"shopApi/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключаемся к базе данных
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	router := gin.Default()

	// Роуты для продуктов
	router.GET("/products", handlers.GetProducts(db))
	router.GET("/products/:id", handlers.GetProduct(db))
	router.POST("/products", handlers.CreateProduct(db))
	router.PUT("/products/:id", handlers.UpdateProduct(db))
	router.DELETE("/products/:id", handlers.DeleteProduct(db))

	// Роуты для корзины
	router.GET("/cart/:userId", handlers.GetCart(db))
	router.POST("/cart/:userId", handlers.AddToCart(db))
	router.DELETE("/cart/:userId/:productId", handlers.RemoveFromCart(db))

	// Роуты для избранного
	router.GET("/favorites/:userId", handlers.GetFavorites(db))
	router.POST("/favorites/:userId", handlers.AddToFavorites(db))
	router.DELETE("/favorites/:userId/:productId", handlers.RemoveFromFavorites(db))

	// Запуск сервера
	router.Run(":8080")
}
