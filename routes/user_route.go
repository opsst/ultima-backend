package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func UserRoute(app *fiber.App) {
	app.Post("/user/login", controllers.Login)
	app.Post("/user/create", controllers.CreateUser)

	app.Post("/user/fb_login", controllers.Fb_Login)
	app.Post("/user/fb_create", controllers.Fb_Create)

	app.Post("/user/google_login", controllers.Google_Login)
	app.Post("/user/google_create", controllers.Google_Create)

	app.Post("/notification/sendall", controllers.PushNotification)

	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte("ultima")}))
	app.Get("/user/ultima", controllers.GetAllUltimaUser)
	app.Get("/user/:userId", controllers.GetAUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Put("/user/point/:userId", controllers.AddUserPoint)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)

}
