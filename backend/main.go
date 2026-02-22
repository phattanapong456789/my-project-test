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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using system env")
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	log.Println("Connected to PostgreSQL")

	db.AutoMigrate(&models.User{}, &models.Table{}, &models.FloorItem{}, &models.Reservation{})
	log.Println("Database migrated")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	auth := handlers.NewAuthHandler(db)
	adminUser := handlers.NewAdminHandler(db)
	reservation := handlers.NewReservationHandler(db)
	tableAdmin := handlers.NewTableAdminHandler(db)

	api := r.Group("/api")
	{
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/login", auth.Login)
		api.GET("/tables", reservation.GetTables) // public: โต๊ะ + floor items + สถานะ

		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", auth.Me)
			protected.PUT("/auth/change-password", auth.ChangePassword)
			protected.POST("/reservations", reservation.CreateReservation)
			protected.GET("/reservations", reservation.GetMyReservations)
			protected.DELETE("/reservations/:id", reservation.CancelReservation)
		}

		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AuthRequired(), middleware.AdminOnly())
		{
			// Users
			adminRoutes.GET("/users", adminUser.GetUsers)
			adminRoutes.POST("/users", adminUser.CreateUser)
			adminRoutes.PUT("/users/:id/role", adminUser.UpdateUserRole)
			adminRoutes.PUT("/users/:id/password", adminUser.UpdateUserPassword)
			adminRoutes.DELETE("/users/:id", adminUser.DeleteUser)

			// Tables
			adminRoutes.GET("/tables", tableAdmin.GetAllTables)
			adminRoutes.POST("/tables", tableAdmin.CreateTable)
			adminRoutes.PUT("/tables/:id", tableAdmin.UpdateTable)
			adminRoutes.DELETE("/tables/:id", tableAdmin.DeleteTable)

			// Floor Items
			adminRoutes.GET("/floor-items", tableAdmin.GetFloorItems)
			adminRoutes.POST("/floor-items", tableAdmin.CreateFloorItem)
			adminRoutes.PUT("/floor-items/:id", tableAdmin.UpdateFloorItem)
			adminRoutes.DELETE("/floor-items/:id", tableAdmin.DeleteFloorItem)

			// Reservations
			adminRoutes.GET("/reservations/summary", tableAdmin.GetSummary)
			adminRoutes.GET("/reservations", tableAdmin.GetAllReservations)
			adminRoutes.PUT("/reservations/:id/status", tableAdmin.UpdateReservationStatus)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on http://localhost:" + port)
	r.Run(":" + port)
}
