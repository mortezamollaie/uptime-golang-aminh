package main

import (
	"github.com/gofiber/fiber/v2"
	"uptime/database"
	"uptime/routes"
)

func main() {
	database.Connect()

	app := fiber.New()
	routes.SetupRoutes(app)

	app.Listen(":3000")
}
