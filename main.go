package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bolean304/e-commerce-cart/controllers"
	"github.com/bolean304/e-commerce-cart/database"
	"github.com/bolean304/e-commerce-cart/middleware"
	"github.com/bolean304/e-commerce-cart/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := controllers.NewApplication(
		database.ProductData(database.Client, database.ProductsCollection),
		database.UserData(database.Client, database.UsersCollection),
	)
	router := gin.New()
	router.Use(gin.Logger())
	// 👇 ADD THIS LINE
	router.Use(middleware.PrometheusMiddleware())
	// Use the CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your React app's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	routes.UserRoutes(router)

	// Protected routes
	auth := router.Group("/")
	auth.Use(middleware.Authentication())

	auth.GET("/addtocart", app.AddToCart())
	auth.GET("/removeitem", app.RemoveItem())
	auth.GET("/instantbuy", app.InstantBuy())
	auth.GET("/checkout", app.BuyFromCart())
	fmt.Printf("port : %v\n", port)
	// Run the server on the specified port
	log.Fatal(router.Run(":" + port))
}
