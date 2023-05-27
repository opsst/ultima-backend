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

var rewardCollection *mongo.Collection = configs.GetCollection2(configs.DB2, "rewards")

func CreateReward(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var reward models.Reward
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&reward); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&reward); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	err := rewardCollection.FindOne(ctx, bson.M{"r_name": reward.R_name}).Decode(&reward)
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "This reward Already have."}})
	}
	newReward := models.Reward{
		ID:          primitive.NewObjectID(),
		R_name:      reward.R_name,
		R_brand:     reward.R_brand,
		R_desc:      reward.R_desc,
		R_point:     reward.R_point,
		R_img:       reward.R_img,
		R_frequency: reward.R_frequency,
	}

	result, err := rewardCollection.InsertOne(ctx, newReward)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllRewards(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	var reward []models.Reward
	defer cancel()

	results, err := rewardCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleReward models.Reward
		if err = results.Decode(&singleReward); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		reward = append(reward, singleReward)
	}

	return c.JSON(fiber.Map{"data": reward})
}

func GetAllRewards_Brand(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	brandName := c.Params("brandName")
	var reward []models.Reward
	defer cancel()

	results, err := rewardCollection.Find(ctx, bson.M{"r_brand": brandName})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleReward models.Reward
		if err = results.Decode(&singleReward); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		reward = append(reward, singleReward)
	}

	return c.JSON(fiber.Map{"data": reward})
}

func DeleteAReward(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	rewardId := c.Params("rewardId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(rewardId)

	result, err := rewardCollection.DeleteOne(ctx, bson.M{"_id": objId})
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

func GetAReward(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var users []models.User
	rewardId := c.Params("rewardId")
	var reward []models.Reward
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(rewardId)
	results, err := rewardCollection.Find(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleReward models.Reward
		if err = results.Decode(&singleReward); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}

		reward = append(reward, singleReward)
	}

	return c.JSON(fiber.Map{"data": reward})

}
