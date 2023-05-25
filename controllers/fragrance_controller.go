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

var fragranceCollection *mongo.Collection = configs.GetCollection2(configs.DB2, "fragrances")

func CreateFragrance(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var fragrance models.Fragrance
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&fragrance); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&fragrance); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	err := ingredientCollection.FindOne(ctx, bson.M{"name": fragrance.P_name}).Decode(&fragrance)
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "This Fragrance Already have."}})
	}
	newFragrance := models.Fragrance{
		ID:      primitive.NewObjectID(),
		P_name:  fragrance.P_name,
		P_brand: fragrance.P_brand,
		P_desc:  fragrance.P_desc,
		P_cate:  fragrance.P_cate,
		P_img:   fragrance.P_img,
		Ing_id:  fragrance.Ing_id,
	}

	result, err := fragranceCollection.InsertOne(ctx, newFragrance)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllFragrances(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var fragrance []models.Fragrance
	defer cancel()

	results, err := fragranceCollection.Find(ctx, bson.M{"p_cate": "Fragrance"})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleFragrance models.Fragrance
		if err = results.Decode(&singleFragrance); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		fragrance = append(fragrance, singleFragrance)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": fragrance}},
	)
}

// GET
// func GetAFragrance(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	fragranceId := c.Params("fragranceId")
// 	var fragrance models.Fragrance
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(fragranceId)

// 	err := fragranceCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&fragrance)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
// 	}

// 	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": fragrance}})
// }

func GetAFragrance(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	fragranceId := c.Params("fragranceId")
	var fragrance []models.Fragrance
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(fragranceId)
	results, err := fragranceCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleFragrance models.Fragrance
		if err = results.Decode(&singleFragrance); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		fragrance = append(fragrance, singleFragrance)
	}

	return c.JSON(fiber.Map{"data": fragrance})

}

// DELETE
func DeleteAFragrance(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fragranceId := c.Params("fragranceId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(fragranceId)

	result, err := fragranceCollection.DeleteOne(ctx, bson.M{"_id": objId})
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

func GetAFragran_ing(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// var users []models.User
	fragranceId := c.Params("fragranceId")
	// var cosmetic []models.Cosmetic

	var myarray []interface{}
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(fragranceId)
	results, err := fragranceCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleFragrance models.Fragrance

		// var ingredient []models.Ingredient
		if err = results.Decode(&singleFragrance); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		for i := 0; i < len(singleFragrance.Ing_id); i++ {
			var ingredient []models.Ingredient
			fmt.Println(singleFragrance.Ing_id[i])
			// fmt.Println(singleFragrance.Cos_name)

			defer cancel()

			results, err := ingredientCollection.Find(ctx, bson.M{"_id": singleFragrance.Ing_id[i]})

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
