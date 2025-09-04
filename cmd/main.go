package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"

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
		log.Println("No .env file found")
	}

	database.Connect()

	go startUptimeChecker()

	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
