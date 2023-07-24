package handlers

import (
	"context"
	"encoding/json"
	"log"

	"example.com/myproject/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

var collection *mongo.Collection

func init() {
	// MongoDB configuration
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("myproject").Collection("users")
}

func GetUsersHandler(c *fiber.Ctx) error {
	// Query all users from MongoDB
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "Failed to get users",
		})
	}
	defer cur.Close(context.Background())

	// Iterate over the cursor and collect users
	users := make(map[string]string)
	for cur.Next(context.Background()) {
		var user User
		if err := cur.Decode(&user); err != nil {
			log.Println(err)
			continue
		}
		users[user.Username] = user.Password
	}

	resp := models.UsersResponse{
		Users: users,
	}
	return c.JSON(resp)
}

func SignupHandler(c *fiber.Ctx) error {
	// Parse request body
	var req models.SignupRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.SignupResponse{
			Message: "Invalid request body",
		})
	}

	// Check if username already exists in MongoDB
	filter := bson.M{"username": req.Username}
	var existingUser User
	if err := collection.FindOne(context.Background(), filter).Decode(&existingUser); err == nil {
		return c.Status(fiber.StatusConflict).JSON(models.SignupResponse{
			Message: "Username already exists",
		})
	}

	// Save new user to MongoDB
	newUser := User{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := collection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "Failed to create user",
		})
	}

	// Assign the inserted ID back to the newUser struct
	newUserID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "Failed to create user",
		})
	}
	newUser.ID = newUserID

	return c.Status(fiber.StatusCreated).JSON(models.SignupResponse{
		Message: "User registered successfully",
	})
}

func LoginHandler(c *fiber.Ctx) error {
	// Parse request body
	var req models.LoginRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.LoginResponse{
			Message: "Invalid request body",
		})
	}

	// Find user in MongoDB
	filter := bson.M{"username": req.Username, "password": req.Password}
	var existingUser User
	if err := collection.FindOne(context.Background(), filter).Decode(&existingUser); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.LoginResponse{
			Message: "Invalid username or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Message: "Login successful",
	})
}
