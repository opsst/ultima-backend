package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	jwtware "github.com/gofiber/jwt/v3"
)

func UserRoute(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))
	app.Post("/user/login", controllers.Login)

	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Post("/user/create", controllers.CreateUser)

	app.Get("/user/:userId", controllers.GetAUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)
}
