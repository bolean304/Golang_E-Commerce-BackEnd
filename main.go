package main

import (
	"log"
	"os"

	"github.com/bolean304/e-commerce-cart/controllers"
	"github.com/bolean304/e-commerce-cart/database"
	"github.com/bolean304/e-commerce-cart/middleware"
	"github.com/bolean304/e-commerce-cart/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/instantbuy", app.InstantBuy())
	router.GET("/checkout", app.BuyFromCart())

	log.Fatal(router.Run(":port"))
}
