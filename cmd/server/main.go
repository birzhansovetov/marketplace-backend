package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Phone string `json:"phone"`
	Items []Item `json:"items,omitempty" gorm:"foreignKey:SellerID"`
}

type Item struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	SellerID    uint    `json:"seller_id"`
	Seller      User    `json:"seller,omitempty"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type CreateItemRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	SellerID    uint    `json:"seller_id" binding:"required"`
}

type UpdateItemRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	SellerID    uint    `json:"seller_id" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active sold archived"`
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=marketplace_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{}, &Item{})
	if err != nil {
		panic("failed to migrate database")
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/users", func(c *gin.Context) {
			var req CreateUserRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user := User{
				Name:  req.Name,
				Email: req.Email,
				Phone: req.Phone,
			}

			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, user)
		})

		api.GET("/users", func(c *gin.Context) {
			var users []User
			if err := db.Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, users)
		})

		api.GET("/users/:id", func(c *gin.Context) {
			var user User
			if err := db.First(&user, c.Param("id")).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusOK, user)
		})

		api.GET("/users/:id/items", func(c *gin.Context) {
			var user User
			if err := db.First(&user, c.Param("id")).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}

			var items []Item
			if err := db.Where("seller_id = ?", c.Param("id")).Find(&items).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user":  user,
				"items": items,
			})
		})

		api.POST("/items", func(c *gin.Context) {
			var req CreateItemRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var seller User
			if err := db.First(&seller, req.SellerID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "seller not found"})
				return
			}

			item := Item{
				Title:       req.Title,
				Description: req.Description,
				Price:       req.Price,
				Status:      "active",
				SellerID:    req.SellerID,
			}

			if err := db.Create(&item).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, item)
		})

		api.GET("/items", func(c *gin.Context) {
			var items []Item
			query := db.Preload("Seller")

			if status := c.Query("status"); status != "" {
				query = query.Where("status = ?", status)
			}

			if err := query.Find(&items).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, items)
		})

		api.GET("/items/:id", func(c *gin.Context) {
			var item Item
			if err := db.Preload("Seller").First(&item, c.Param("id")).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		api.PUT("/items/:id", func(c *gin.Context) {
			var item Item
			if err := db.First(&item, c.Param("id")).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
				return
			}

			var req UpdateItemRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			item.Title = req.Title
			item.Description = req.Description
			item.Price = req.Price
			item.SellerID = req.SellerID

			if err := db.Save(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, item)
		})

		api.DELETE("/items/:id", func(c *gin.Context) {
			result := db.Delete(&Item{}, c.Param("id"))
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			if result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
		})

		api.PATCH("/items/:id/status", func(c *gin.Context) {
			var item Item
			if err := db.First(&item, c.Param("id")).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
				return
			}

			var req UpdateStatusRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			item.Status = req.Status
			if err := db.Save(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, item)
		})
	}

	r.Run(":8080")
}
