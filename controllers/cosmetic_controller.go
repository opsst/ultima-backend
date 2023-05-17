package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cosmeticCollection *mongo.Collection = configs.GetCollection2(configs.DB2, "cosmetics")

func CreateCosmetic(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var cosmetic models.Cosmetic
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&cosmetic); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&cosmetic); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newCosmetic := models.Cosmetic{
		Id:      primitive.NewObjectID(),
		P_name:  cosmetic.P_name,
		P_brand: cosmetic.P_brand,
		P_desc:  cosmetic.P_desc,
		P_cate:  cosmetic.P_cate,
		P_img:   cosmetic.P_img,
		Ing_id:  cosmetic.Ing_id,
	}

	result, err := cosmeticCollection.InsertOne(ctx, newCosmetic)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}
