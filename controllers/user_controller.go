package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/lestrrat-go/jwx/v2/jwa"
	jwt2 "github.com/lestrrat-go/jwx/v2/jwt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid body."})
	}

	if user.Email == "" {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Email could not be empty."})
	}

	if user.Password == "" {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Password could not be empty."})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	if !strings.Contains(user.Email, "@") {
		fmt.Println(user.Email)
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Wrong email format."})
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
	if err == nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "This email already taken."})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	newUser := models.User{
		Id:             primitive.NewObjectID(),
		Email:          user.Email,
		Password:       string(hash),
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		Admin:          "NA",
		Used_Point_URL: user.Used_Point_URL,
		Point:          0,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid data."})
	}

	claims := jwt.MapClaims{
		"id":     newUser.Id,
		"email":  newUser.Email,
		"admin":  newUser.Admin,
		"f_name": newUser.Firstname,
		"l_name": newUser.Lastname,
		// "exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ultima"))
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Fail to get token"})
	}
	return c.JSON(fiber.Map{"status": http.StatusOK, "message": "Success", "result": result, "token": t})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid body."})
	}

	if !strings.Contains(user.Email, "@") {
		fmt.Println(user.Email)
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Wrong email format."})
	}

	if user.Email == "" {
		// fmt.Println(user.Email)
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Email could not be empty."})
	}

	if user.Password == "" {
		// fmt.Println(user.Email)
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Password could not be empty."})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}
	var plainpassword = user.Password
	var firebase_token = user.Firebasetoken
	// , "password": user.Password
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Could not find your email."})
	}

	// println(user.Password)
	byteHash := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(plainpassword))
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Your password is wrong."})
	}
	// fmt.Println("kuy")
	// fmt.Println(firebase_token)
	update := bson.M{"firebasetoken": firebase_token}

	result, err := userCollection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": update})
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Can not update firebase token."})
	}

	claims := jwt.MapClaims{
		"id":     user.Id,
		"email":  user.Email,
		"admin":  user.Admin,
		"f_name": user.Firstname,
		"l_name": user.Lastname,
		// "exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ultima"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// fmt.Println(result, t)
	return c.JSON(fiber.Map{"token": t, "message": "success", "result": result, "status": http.StatusOK})
	// return fiber.ErrServiceUnavailable
}

func Fb_Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid body."})
	}

	if user.Fb_login == "" {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Someting might wrong."})
	}

	err := userCollection.FindOne(ctx, bson.M{"fb_login": user.Fb_login}).Decode(&user)
	if err == nil {
		claims := jwt.MapClaims{
			"id":       user.Id,
			"fb_login": user.Fb_login,
			// "exp":   time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("ultima"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"token": t, "message": "success with used signup uid", "status": http.StatusOK})

	}
	newUser := models.User{
		Id:             primitive.NewObjectID(),
		Fb_login:       user.Fb_login,
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		Admin:          "NA",
		Used_Point_URL: user.Used_Point_URL,
		Point:          0,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid data."})
	}

	claims := jwt.MapClaims{
		"id":       newUser.Id,
		"fb_login": user.Fb_login,
		// "exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ultima"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t, "message": "success", "result": result, "status": http.StatusOK})
}

func Fb_Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid body."})
	}

	if user.Fb_login == "" {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Someting might wrong."})
	}

	var firebase_token = user.Firebasetoken
	// , "password": user.Password
	err := userCollection.FindOne(ctx, bson.M{"fb_login": user.Fb_login}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Could not find your uid."})
	}

	// fmt.Println("kuy")
	// fmt.Println(firebase_token)
	update := bson.M{"firebasetoken": firebase_token}

	result, err := userCollection.UpdateOne(ctx, bson.M{"fb_login": user.Fb_login}, bson.M{"$set": update})
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Can not update firebase token."})
	}

	claims := jwt.MapClaims{
		"id":       user.Id,
		"fb_login": user.Fb_login,
		// "exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ultima"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t, "message": "success", "result": result, "status": http.StatusOK})
}

func Google_Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid body."})
	}

	if user.Google_login == "" {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Someting might wrong with your facebook account."})
	}

	err := userCollection.FindOne(ctx, bson.M{"google_login": user.Google_login}).Decode(&user)
	if err == nil {
		claims := jwt.MapClaims{
			"id":           user.Id,
			"google_login": user.Google_login,
			// "exp":   time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("ultima"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"token": t, "message": "success with used signup uid", "status": http.StatusOK})

	}
	newUser := models.User{
		Id:             primitive.NewObjectID(),
		Google_login:   user.Google_login,
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		Admin:          "NA",
		Used_Point_URL: user.Used_Point_URL,
		Point:          0,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid data."})
	}

	claims := jwt.MapClaims{
		"id":           newUser.Id,
		"google_login": user.Google_login,
		// "exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ultima"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t, "message": "success", "result": result, "status": http.StatusOK})
}

func Google_Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Invalid body."})
	}

	if user.Google_login == "" {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Someting might wrong with your gmail."})
	}

	var firebase_token = user.Firebasetoken
	// , "password": user.Password
	err := userCollection.FindOne(ctx, bson.M{"google_login": user.Google_login}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Could not find your uid."})
	}

	// fmt.Println("kuy")
	// fmt.Println(firebase_token)
	update := bson.M{"firebasetoken": firebase_token}

	result, err := userCollection.UpdateOne(ctx, bson.M{"google_login": user.Google_login}, bson.M{"$set": update})
	if err != nil {
		return c.JSON(fiber.Map{"status": http.StatusInternalServerError, "message": "Can not update firebase token."})
	}

	claims := jwt.MapClaims{
		"id":       user.Id,
		"fb_login": user.Fb_login,
		// "exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ultima"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t, "message": "success", "result": result, "status": http.StatusOK})
}

func GetATokenDetail(c *fiber.Ctx) error {
	myUserId := strings.Split(c.Get("Authorization"), " ")
	var token []string

	tok, err := jwt2.Parse([]byte(myUserId[1]), jwt2.WithKey(jwa.HS256, []byte("ultima")))
	if err != nil {
		fmt.Println(token[1])
	}
	return c.JSON(fiber.Map{"data": tok})
}
func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()
	userId := c.Params("userId")
	var fb = false
	var google = false
	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	user.Password = ""
	if user.Fb_login != "" {
		user.Fb_login = ""
		fb = true
	}
	if user.Google_login != "" {
		user.Google_login = ""
		google = true
	}

	return c.JSON(fiber.Map{"data": user, "fb": fb, "google": google, "point": user.Point, "status": http.StatusOK})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// update := bson.M{"email": user.Email, "password": user.Password, "firstname": user.Firstname, "lastname": user.Lastname, "admin": user.Admin}
	update := bson.M{"admin": user.Admin}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
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

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	// return c.Status(http.StatusOK).JSON(
	// 	responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users ,"status": http.StatusOK}},
	// )
	return c.JSON(fiber.Map{"data": users, "status": http.StatusOK})
}

func GetAllUltimaUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	filter := bson.D{{"email", primitive.Regex{Pattern: "@ultima.com", Options: ""}}}
	results, err := userCollection.Find(ctx, filter)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

func AddUserPoint(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	// var users []models.User
	var user models.User
	defer cancel()
	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	objId, _ := primitive.ObjectIDFromHex(userId)
	results, err := userCollection.Find(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"Error on result: ": err.Error()}})
		}
		// users = append(users, singleUser)
		for i := 0; i < len(singleUser.Used_Point_URL); i++ {
			if user.Used_Point_URL[0] == singleUser.Used_Point_URL[i] {
				return c.JSON(fiber.Map{"message": "Used Link", "status": http.StatusInternalServerError})
			}

		}
		user.Point = singleUser.Point + 1000
		user.Used_Point_URL = append(singleUser.Used_Point_URL, user.Used_Point_URL...)
		fmt.Println(user.Used_Point_URL...)
		// fmt.Println(singleUser.Used_Point_URL)
	}
	// return c.JSON(fiber.Map{"message": users})

	// update := bson.M{"email": user.Email, "password": user.Password, "firstname": user.Firstname, "lastname": user.Lastname, "admin": user.Admin}
	update := bson.M{"point": user.Point, "used_point_url": user.Used_Point_URL}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
	// return c.JSON(fiber.Map{"data": user, "status": http.StatusOK})
}

func DelUserPoint(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	// var users []models.User
	var user models.User
	defer cancel()
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var redeemPoint = user.Point

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	fmt.Println(user.Point)
	fmt.Println(redeemPoint)
	var totalPoint = user.Point - redeemPoint
	// update := bson.M{"email": user.Email, "password": user.Password, "firstname": user.Firstname, "lastname": user.Lastname, "admin": user.Admin}
	update := bson.M{"point": totalPoint}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
	// return c.JSON(fiber.Map{"data": user, "status": http.StatusOK})
}

func PushNotification(c *fiber.Ctx) error {
	apps, _, _ := configs.SetupFirebase()
	// fmt.Println(apps)

	type Notification struct {
		Title string
		Body  string
	}
	var noti Notification

	if err := c.BodyParser(&noti); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	var title = noti.Title
	var body = noti.Body
	fmt.Println(title + " : " + body)
	sendToToken(apps, title, body)
	return c.JSON(fiber.Map{"message": "yey", "status": http.StatusOK})

}

func sendToToken(apps *firebase.App, title string, body string) {
	ctx := context.Background()
	client, err := apps.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v", err)
	}

	registrationToken := "fU_KNvolTA2FCgjO7L15K6:APA91bGCkY5iyiNj4cE0S3nh05PKyXwombEJ__PcuJh-bOSLSWgz_XNrt50g6u4-cMEXnVw90y6svez-9DxNW8gopj0sfecSCQvcNo4cBqCWfGB6HKYpewcA4_zXlo6x4zP-MJdoIiHo"

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: registrationToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent message:", response)
}
