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

var skincareCollection *mongo.Collection = configs.GetCollection2(configs.DB2, "skincares")

func CreateSkincare(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var skincare models.Skincare
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&skincare); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&skincare); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	err := skincareCollection.FindOne(ctx, bson.M{"p_name": skincare.P_name}).Decode(&skincare)
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "This Skincare Already have."}})
	}
	newSkincare := models.Skincare{
		ID:      primitive.NewObjectID(),
		P_name:  skincare.P_name,
		P_brand: skincare.P_brand,
		P_desc:  skincare.P_desc,
		P_cate:  skincare.P_cate,
		P_img:   skincare.P_img,
		L_link:  skincare.L_link,
		Ing_id:  skincare.Ing_id,
	}

	result, err := skincareCollection.InsertOne(ctx, newSkincare)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllSkincares(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var skincare []models.Skincare
	defer cancel()

	results, err := skincareCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSkincare models.Skincare
		if err = results.Decode(&singleSkincare); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		skincare = append(skincare, singleSkincare)
	}

	return c.JSON(fiber.Map{"data": skincare})
}

// GET
// func GetASkincare(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	skincareId := c.Params("skincareId")
// 	var skincare models.Skincare
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(skincareId)

// 	err := skincareCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&skincare)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
// 	}

// 	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": skincare}})
// }

func GetASkincare(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	skincareId := c.Params("skincareId")
	var skincare []models.Skincare
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(skincareId)
	results, err := skincareCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSkincare models.Skincare
		if err = results.Decode(&singleSkincare); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		skincare = append(skincare, singleSkincare)
	}

	return c.JSON(fiber.Map{"data": skincare})

}

// PUT
func EditASkincare(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	skincareId := c.Params("skincareId")
	var skincare models.Skincare
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(skincareId)

	//validate the request body
	if err := c.BodyParser(&skincare); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&skincare); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"p_name": skincare.P_name, "p_brand": skincare.P_brand}

	result, err := skincareCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedSkincare models.Skincare
	if result.MatchedCount == 1 {
		err := skincareCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedSkincare)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedSkincare}})
}

// DELETE
func DeleteASkincare(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	skincareId := c.Params("skincareId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(skincareId)

	result, err := skincareCollection.DeleteOne(ctx, bson.M{"_id": objId})
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

func GetASkincare_ing(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// var users []models.User
	skincareId := c.Params("skincareId")
	// var cosmetic []models.Cosmetic

	var myarray []interface{}
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(skincareId)
	results, err := skincareCollection.Find(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSkincareId models.Skincare

		// var ingredient []models.Ingredient
		if err = results.Decode(&singleSkincareId); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		for i := 0; i < len(singleSkincareId.Ing_id); i++ {
			var ingredient []models.Ingredient
			// fmt.Println(singleSkincareId.Ing_id[i])
			// fmt.Println(singleSkincareId.Cos_name)

			defer cancel()

			results, err := ingredientCollection.Find(ctx, bson.M{"_id": singleSkincareId.Ing_id[i]})

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
