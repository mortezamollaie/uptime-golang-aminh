package controllers

import (
	"time"
	"uptime/database"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck checks the health status of the application
// @Summary Health check
// @Description Check the health status of the API and database connection
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Failure 503 {object} map[string]interface{} "Service is unhealthy"
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	// Check database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":    "unhealthy",
			"database":  "disconnected",
			"timestamp": time.Now().Unix(),
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":    "unhealthy",
			"database":  "unreachable",
			"timestamp": time.Now().Unix(),
		})
	}

	return c.JSON(fiber.Map{
		"status":    "healthy",
		"database":  "connected",
		"timestamp": time.Now().Unix(),
	})
}
