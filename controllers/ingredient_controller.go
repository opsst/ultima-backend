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

var ingredientCollection *mongo.Collection = configs.GetCollection2(configs.DB2, "cosmetic-ingredients")

func CreateIngredient(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var ingredient models.Ingredient
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&ingredient); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&ingredient); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	err := ingredientCollection.FindOne(ctx, bson.M{"name": ingredient.Name}).Decode(&ingredient)
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "This Ingredient Already have."}})
	}

	newIngredient := models.Ingredient{
		Name:    ingredient.Name,
		Rate:    ingredient.Rate,
		Calling: ingredient.Calling,
		Func:    ingredient.Func,
		Irr:     ingredient.Irr,
		Come:    ingredient.Come,
		Cosing:  ingredient.Cosing,
		Quick:   ingredient.Quick,
		Detail:  ingredient.Detail,
		Proof:   ingredient.Proof,
		Link:    ingredient.Link,
		IsTryOn: ingredient.IsTryOn,
	}

	result, err := ingredientCollection.InsertOne(ctx, newIngredient)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.JSON(fiber.Map{"ing_id": result.InsertedID})
}

func GetAllIngredients(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var ingredient []models.Ingredient
	defer cancel()

	results, err := ingredientCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleIngredient models.Ingredient
		if err = results.Decode(&singleIngredient); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		ingredient = append(ingredient, singleIngredient)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": ingredient}},
	)
}
