package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bolean304/e-commerce-cart/controllers"
	"github.com/bolean304/e-commerce-cart/database"
	"github.com/bolean304/e-commerce-cart/middleware"
	"github.com/bolean304/e-commerce-cart/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(
		database.ProductData(database.Client, "Products"),
		database.UserData(database.Client, "Users"),
	)
	router := gin.New()
	router.Use(gin.Logger())
	// Use the CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your React app's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/instantbuy", app.InstantBuy())
	router.GET("/checkout", app.BuyFromCart())
	fmt.Println("port : %v", port)
	// Run the server on the specified port
	log.Fatal(router.Run(":" + port))
}
