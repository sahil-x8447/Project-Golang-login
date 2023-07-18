package main

import (
	"log"

	"example.com/myproject/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Routes
	app.Get("/users", handlers.GetUsersHandler)
	app.Post("/signup", handlers.SignupHandler)
	app.Post("/login", handlers.LoginHandler)

	// Start server
	err := app.Listen(":6000")
	if err != nil {
		log.Fatal(err)
	}
}
