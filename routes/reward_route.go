package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
)

func RewardRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Post("/reward/create", controllers.CreateReward)
	// app.Get("/skincare/checkall", controllers.GetAllSkincares)
	// app.Get("/skincare/ingredient/:skincareId", controllers.GetASkincare_ing)
	app.Get("/reward/checkall", controllers.GetAllRewards)
	app.Get("/reward/brand/:brandName", controllers.GetAllRewards_Brand)
	app.Get("/reward/:rewardId", controllers.GetAReward)
	// app.Put("/skincare/:skincareId", controllers.EditASkincare)
	app.Delete("/reward/:rewardId", controllers.DeleteAReward)

}
