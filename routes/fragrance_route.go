package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
)

func FragranceRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Post("/fragrance/create", controllers.CreateFragrance)
	app.Get("/fragrance/checkall", controllers.GetAllFragrances)
	app.Get("/fragrance/ingredient/:fragranceId", controllers.GetAFragran_ing)
	// app.Get("/fragrance/getfragrance", controllers.GetAllFragances)
	app.Get("/fragrance/:fragranceId", controllers.GetAFragrance)
	app.Put("/fragrance/:fragranceId", controllers.EditASkincare)
	app.Delete("/fragrance/:fragranceId", controllers.DeleteAFragrance)

}
