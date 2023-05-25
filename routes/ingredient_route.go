package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	jwtware "github.com/gofiber/jwt/v3"
)

func IngredientRoute(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))

	app.Post("/ingredient/find", controllers.FindIngredient)
	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Post("/ingredient/create", controllers.CreateIngredient)
	app.Get("/ingredient/checkall", controllers.GetAllIngredients)
}
