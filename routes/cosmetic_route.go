package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// jwtware "github.com/gofiber/jwt/v3"
)

func CosmeticRoute(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))
	app.Get("/cosmetic/checktryon", controllers.GetAllTryonCosmetics)
	app.Get("/cosmetic/checkall", controllers.GetAllCosmetics)
	app.Get("/cosmetic/ingredient/:cosId", controllers.GetACosmetic_ing)

	// app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Post("/cosmetic/create", controllers.CreateCosmetic)
	app.Delete("/cosmetic/:cosId", controllers.DeleteACosmetic)

}
