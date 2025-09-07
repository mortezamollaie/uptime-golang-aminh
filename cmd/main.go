package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/robfig/cron/v3"

	"uptime/config"
	"uptime/database"
	"uptime/models"
	"uptime/routes"
	"uptime/uptime"
)

func startUptimeChecker() *cron.Cron {
	c := cron.New()

	checkInterval := "@every " + config.AppConfig.UptimeChecker.CheckInterval.String()
	c.AddFunc(checkInterval, func() {
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
	return c
}

func main() {
	config.Load()
	
	database.Connect()

	cron := startUptimeChecker()

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

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetupRoutes(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		port := config.AppConfig.Server.Port
		log.Printf("Server starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	
	cron.Stop()
	
	if err := app.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
