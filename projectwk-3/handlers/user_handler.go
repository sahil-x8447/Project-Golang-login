package handlers

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"example.com/myproject/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var jwtSecretKey = []byte("secret-key")

const JWTTokenDuration = time.Hour * 24 // Token expires in 24 hours

// GenerateJWT generates a new JWT token with the provided username.
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(JWTTokenDuration)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = expirationTime.Unix()

	// Add any additional claims if needed, e.g., role, expiration, etc.

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT validates the JWT token provided in the Authorization header.
// If the token is valid, it returns the username extracted from the token.
func ValidateJWT(c *fiber.Ctx) (string, error) {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if tokenString == "" {
		return "", errors.New("missing JWT token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, found := claims["username"].(string)
		if !found {
			return "", errors.New("invalid token claims")
		}
		return username, nil
	}

	return "", errors.New("invalid JWT token")
}

var collection *mongo.Collection

func init() {
	// MongoDB configuration
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	collection = client.Database("myproject").Collection("users")
}

// GetUserDetailsHandler retrieves details of the logged-in user based on the provided username.
func GetUserDetailsHandler(req models.GetUserDetailsRequest) (models.GetUserDetailsResponse, error) {
	// Find user in MongoDB by username
	filter := bson.M{"username": req.Username}
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return models.GetUserDetailsResponse{}, err
	}

	return models.GetUserDetailsResponse{
		Username: user.Username,
		// You can include any additional user details you want to return here
	}, nil
}

func SignupHandler(req models.SignupRequest) (models.SignupResponse, error) {
	// Check if username already exists in MongoDB
	filter := bson.M{"username": req.Username}
	var existingUser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.FindOne(ctx, filter).Decode(&existingUser); err == nil {
		return models.SignupResponse{
			Message: "Username already exists",
		}, nil
	}

	// Save new user to MongoDB
	newUser := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := collection.InsertOne(ctx, newUser); err != nil {
		return models.SignupResponse{}, err
	}

	return models.SignupResponse{
		Message: "User registered successfully",
	}, nil
}

func LoginHandler(req models.LoginRequest) (models.LoginResponse, error) {
	// Find user in MongoDB
	filter := bson.M{"username": req.Username, "password": req.Password}
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return models.LoginResponse{
			Message: "Invalid username or password",
		}, nil
	}

	// Generate JWT token
	token, err := GenerateJWT(req.Username)
	if err != nil {
		log.Println("Failed to generate JWT token:", err)
		return models.LoginResponse{
			Message: "Login successful, but failed to generate token",
		}, nil
	}

	return models.LoginResponse{
		Message: "Login successful",
		Token:   token,
	}, nil
}
