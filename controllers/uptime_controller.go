package controllers

import (
	"uptime/database"
	"uptime/models"
	"uptime/monitoring"

	"github.com/gofiber/fiber/v2"
)

func CheckUptime(c *fiber.Ctx) error {
	var nodes []models.Node
	if err := database.DB.Find(&nodes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if len(nodes) == 0 {
		return c.JSON(fiber.Map{"message": "No nodes found"})
	}

	monitoring.Check(nodes)

	return c.JSON(fiber.Map{"message": "Uptime check completed"})
}
