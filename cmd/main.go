package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"uptime/database"
	"uptime/models"
	"uptime/routes"
	"uptime/uptime"
)

func startUptimeChecker() {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		var nodes []models.Node
		if err := database.DB.Find(&nodes).Error; err != nil {
			log.Println("Error fetching nodes:", err)
			return
		}
		if len(nodes) == 0 {
			log.Println("No nodes found")
			return
		}

		uptime.Check(nodes)
		log.Println("Uptime check completed")
	})

	c.Start()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	database.Connect()

	go startUptimeChecker()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
