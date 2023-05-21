package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

	err := cosmeticCollection.FindOne(ctx, bson.M{"name": cosmetic.Cos_name}).Decode(&cosmetic)
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "This Cosmetic Already have in database."}})
	}
	newCosmetic := models.Cosmetic{
		Cos_brand:       cosmetic.Cos_brand,
		Cos_name:        cosmetic.Cos_name,
		Cos_desc:        cosmetic.Cos_desc,
		Cos_cate:        cosmetic.Cos_cate,
		Cos_img:         cosmetic.Cos_img,
		Cos_istryon:     cosmetic.Cos_istryon,
		Cos_color_img:   cosmetic.Cos_color_img,
		Cos_tryon_name:  cosmetic.Cos_tryon_name,
		Cos_tryon_color: cosmetic.Cos_tryon_color,
		Ing_id:          cosmetic.Ing_id,
	}

	result, err := cosmeticCollection.InsertOne(ctx, newCosmetic)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllCosmetics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var cosmetic []models.Cosmetic
	defer cancel()

	results, err := cosmeticCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleCosmetic models.Cosmetic
		if err = results.Decode(&singleCosmetic); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		cosmetic = append(cosmetic, singleCosmetic)
	}

	return c.JSON(fiber.Map{"data": cosmetic})

}
