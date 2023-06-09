package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ingredientCollection *mongo.Collection = configs.GetCollection2(configs.DB2, "ingredients")

func FindIngredient(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var ingredient models.Ingredient
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&ingredient); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	fmt.Println(ingredient.Name)

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&ingredient); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}
	query := bson.M{
		"$or": bson.A{
			bson.M{"calling": bson.M{
				"$regex": primitive.Regex{Pattern: ingredient.Name, Options: "i"}}},
			bson.M{"name": bson.M{
				"$regex": primitive.Regex{Pattern: ingredient.Name, Options: "i"}}},
		},
	}

	err := ingredientCollection.FindOne(ctx, query).Decode(&ingredient)
	fmt.Println(ingredient.Calling)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "Not Found"})
	}

	return c.JSON(fiber.Map{"ing_id": ingredient.ID.Hex()})

}

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
		return c.JSON(fiber.Map{"ing_id": ingredient.ID})
	}

	newIngredient := models.Ingredient{
		ID:      primitive.NewObjectID(),
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

func GetAIngredient(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// var users []models.User
	ingredientId := c.Params("ingredientId")
	var ingredient []models.Ingredient

	// var myarray []interface{}
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(ingredientId)
	results, err := ingredientCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleIngredient models.Ingredient

		// var ingredient []models.Ingredient
		if err = results.Decode(&singleIngredient); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		ingredient = append(ingredient, singleIngredient)
	}
	// fmt.Println(ingredient)
	return c.JSON(fiber.Map{"data": ingredient})

}
