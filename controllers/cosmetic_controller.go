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

func GetAllCosmetics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var cosmetics []models.Cosmetic
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

		cosmetics = append(cosmetics, singleCosmetic)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": cosmetics}},
	)
}

// GET
func GetACosmetic(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cosmeticId := c.Params("cosmeticId")
	var cosmetic models.Cosmetic
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(cosmeticId)

	err := cosmeticCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&cosmetic)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": cosmetic}})
}

// PUT
func EditACosmetic(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cosmeticId := c.Params("cosmeticId")
	var cosmetic models.Cosmetic
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(cosmeticId)

	//validate the request body
	if err := c.BodyParser(&cosmetic); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&cosmetic); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"p_name": cosmetic.P_name, "p_brand": cosmetic.P_brand}

	result, err := cosmeticCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedCosmetic models.Cosmetic
	if result.MatchedCount == 1 {
		err := cosmeticCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedCosmetic)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedCosmetic}})
}

// DELETE
func DeleteACosmetic(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cosmeticId := c.Params("cosmeticId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(cosmeticId)

	result, err := cosmeticCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}
