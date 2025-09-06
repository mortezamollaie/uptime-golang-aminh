package controllers

import (
	"strconv"
	"strings"
	"uptime/database"
	"uptime/models"
	"uptime/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateNodeLog(c *fiber.Ctx) error {
	var log models.NodeLog
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// Validate NodeID
	if log.NodeID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "NodeID is required"})
	}

	newLog, err := services.CreateNodeLog(&log)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create node log"})
	}
	return c.Status(201).JSON(newLog)
}

func GetAllNodeLogs(c *fiber.Ctx) error {
	logs, err := services.GetAllNodeLogs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}

func GetNodeLog(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	log, err := services.GetNodeLog(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Node log not found"})
	}
	return c.JSON(log)
}

func UpdateNodeLog(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	var body models.NodeLog
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	log, err := services.GetNodeLog(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Node log not found"})
	}

	log.Delay = body.Delay
	log.Status = body.Status
	log.Up = body.Up
	log.Suspended = body.Suspended
	log.Exception = body.Exception

	updatedLog, err := services.UpdateNodeLog(log)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update node log"})
	}
	return c.JSON(updatedLog)
}

func DeleteNodeLog(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	err = services.DeleteNodeLogByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(fiber.Map{"error": "Node log not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete node log"})
	}
	return c.SendStatus(204)
}

func GetAllNodesWithLogs(c *fiber.Ctx) error {
	var nodes []models.Node

	err := database.DB.Preload("NodeLogs", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, node_id, delay, status, up, suspended, exception, created_at")
	}).Find(&nodes).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(nodes)
}
