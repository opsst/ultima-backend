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
		Id:              primitive.NewObjectID(),
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
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		cosmetic = append(cosmetic, singleCosmetic)
	}

	return c.JSON(fiber.Map{"data": cosmetic})
	// return c.JSON(fiber.Map{"data": "hi"})
}

func GetAllTryonCosmetics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var cosmetic []models.Cosmetic
	defer cancel()

	results, err := cosmeticCollection.Find(ctx, bson.M{"cos_istryon": true})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleCosmetic models.Cosmetic
		if err = results.Decode(&singleCosmetic); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		cosmetic = append(cosmetic, singleCosmetic)
	}

	return c.JSON(fiber.Map{"data": cosmetic})

}

func GetACosmetic(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// var users []models.User
	cosId := c.Params("cosId")
	var cosmetic []models.Cosmetic

	// var myarray []interface{}
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(cosId)
	results, err := cosmeticCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleCosmetic models.Cosmetic

		// var ingredient []models.Ingredient
		if err = results.Decode(&singleCosmetic); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		cosmetic = append(cosmetic, singleCosmetic)
	}
	// fmt.Println(ingredient)
	return c.JSON(fiber.Map{"data": cosmetic})

}

func GetACosmetic_ing(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// var users []models.User
	cosId := c.Params("cosId")
	// var cosmetic []models.Cosmetic

	var myarray []interface{}
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(cosId)
	results, err := cosmeticCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleCosmetic models.Cosmetic

		// var ingredient []models.Ingredient
		if err = results.Decode(&singleCosmetic); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		for i := 0; i < len(singleCosmetic.Ing_id); i++ {
			var ingredient []models.Ingredient
			fmt.Println(singleCosmetic.Ing_id[i])
			// fmt.Println(singleCosmetic.Cos_name)

			defer cancel()

			results, err := ingredientCollection.Find(ctx, bson.M{"_id": singleCosmetic.Ing_id[i]})

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
				// fmt.Println("Before Result!")
				ingredient = append(ingredient, singleIngredient)

			}
			myarray = append(myarray, ingredient)
		}
		// cosmetic = append(cosmetic, singleCosmetic)
	}
	// fmt.Println(ingredient)
	return c.JSON(fiber.Map{"data": myarray})

}

func DeleteACosmetic(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cosId := c.Params("cosId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(cosId)

	result, err := cosmeticCollection.DeleteOne(ctx, bson.M{"_id": objId})
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
