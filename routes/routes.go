package routes

import (
	"marketplace-api/handlers"
	"marketplace-api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	userHandler := handlers.NewUserHandler(db)
	itemHandler := handlers.NewItemHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	api.POST("/users", userHandler.CreateUser)
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.GET("/users/:id/items", userHandler.GetUserItems)

	api.GET("/items", itemHandler.GetItems)
	api.GET("/items/:id", itemHandler.GetItemByID)
	api.PUT("/items/:id", itemHandler.UpdateItem)
	api.DELETE("/items/:id", itemHandler.DeleteItem)
	api.PATCH("/items/:id/status", itemHandler.UpdateItemStatus)

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	auth := api.Group("/", middleware.AuthMiddleware())
	{
		auth.POST("/items", itemHandler.CreateItem)

		auth.POST("/orders", orderHandler.CreateOrder)
		auth.GET("/orders", orderHandler.GetOrders)
		auth.GET("/orders/:id", orderHandler.GetOrderByID)
		auth.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
		auth.DELETE("/orders/:id", orderHandler.DeleteOrder)
	}
}
