package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes" //add this

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// sendToToken(apps)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	//run database
	configs.ConnectDB()
	//routes

	routes.CosmeticRoute(app)
	routes.UserRoute(app) //add this
	routes.FragranceRoute(app)
	routes.IngredientRoute(app)
	routes.SkincareRoute(app)
	app.Listen(":8000")
}
