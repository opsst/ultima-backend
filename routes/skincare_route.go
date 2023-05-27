package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
)

func SkincareRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Post("/skincare/create", controllers.CreateSkincare)
	app.Get("/skincare/checkall", controllers.GetAllSkincares)
	// app.Get("/skincare/getfragrance", controllers.GetAllFragances)
	app.Get("/skincare/ingredient/:skincareId", controllers.GetASkincare_ing)
	app.Get("/skincare/:skincareId", controllers.GetASkincare)
	app.Put("/skincare/:skincareId", controllers.EditASkincare)
	app.Delete("/skincare/:skincareId", controllers.DeleteASkincare)

}
