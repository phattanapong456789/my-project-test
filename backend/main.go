package main

import (
	"auth-app/database"
	"auth-app/handlers"
	"auth-app/middleware"
	"auth-app/models"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// โหลด .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using system env")
	}

	// เชื่อม Database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	log.Println("Connected to PostgreSQL")

	// Auto-migrate สร้างตาราง users อัตโนมัติ
	db.AutoMigrate(&models.User{})
	log.Println("Database migrated")

	// Setup Router
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	auth := handlers.NewAuthHandler(db)
	admin := handlers.NewAdminHandler(db)

	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/login", auth.Login)

		// Protected routes (ต้อง login)
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", auth.Me)
			protected.PUT("/auth/change-password", auth.ChangePassword)
		}

		// Admin routes (ต้องเป็น admin เท่านั้น)
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AuthRequired(), middleware.AdminOnly())
		{
			adminRoutes.GET("/users", admin.GetUsers)
			adminRoutes.POST("/users", admin.CreateUser)
			adminRoutes.PUT("/users/:id/role", admin.UpdateUserRole)
			adminRoutes.DELETE("/users/:id", admin.DeleteUser)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on http://localhost:" + port)
	r.Run(":" + port)
}
