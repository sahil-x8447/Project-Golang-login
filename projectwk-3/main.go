package main

import (
	"encoding/json"
	"log"

	"example.com/myproject/handlers"
	"example.com/myproject/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Routes
	app.Get("/users", func(c *fiber.Ctx) error {
		var req models.GetUserDetailsRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Message: "Invalid request body",
			})
		}

		resp, err := handlers.GetUserDetailsHandler(req)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Message: "Failed to get user details",
			})
		}
		return c.JSON(resp)
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		var req models.SignupRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(models.SignupResponse{
				Message: "Invalid request body",
			})
		}

		resp, err := handlers.SignupHandler(req)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Message: "Failed to create user",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(resp)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		var req models.LoginRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(models.LoginResponse{
				Message: "Invalid request body",
			})
		}

		resp, err := handlers.LoginHandler(req)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusUnauthorized).JSON(models.LoginResponse{
				Message: "Invalid username or password",
			})
		}
		return c.JSON(resp)
	})

	// Start server
	err := app.Listen(":6000")
	if err != nil {
		log.Fatal(err)
	}
}
