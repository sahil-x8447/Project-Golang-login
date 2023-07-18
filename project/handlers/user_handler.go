package handlers

import (
	"encoding/json"
	"log"

	"example.com/myproject/models"
	"github.com/gofiber/fiber/v2"
)

var users = make(map[string]string)

func GetUsersHandler(c *fiber.Ctx) error {
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

	// Check if username already exists
	if _, exists := users[req.Username]; exists {
		return c.Status(fiber.StatusConflict).JSON(models.SignupResponse{
			Message: "Username already exists",
		})
	}

	// Save new user
	users[req.Username] = req.Password

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

	// Check if user exists
	password, exists := users[req.Username]
	if !exists || password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(models.LoginResponse{
			Message: "Invalid username or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Message: "Login successful",
	})
}
