package controllers

import (
	"strconv"
	"strings"
	"uptime/models"
	"uptime/services"

	"github.com/gofiber/fiber/v2"
)

func CreateHistory(c *fiber.Ctx) error {
	var history models.History
	if err := c.BodyParser(&history); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// Validate NodeID
	if history.NodeID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "NodeID is required"})
	}

	newHistory, err := services.CreateHistory(&history)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create history"})
	}
	return c.Status(201).JSON(newHistory)
}

func GetAllHistories(c *fiber.Ctx) error {
	histories, err := services.GetAllHistories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(histories)
}

func GetHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	history, err := services.GetHistory(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "History not found"})
	}
	return c.JSON(history)
}

func UpdateHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	var body models.History
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	history, err := services.GetHistory(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "History not found"})
	}

	// Update fields
	history.Delay = body.Delay
	history.Status = body.Status
	history.Up = body.Up
	history.Suspended = body.Suspended
	history.Exception = body.Exception

	updatedHistory, err := services.UpdateHistory(history)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update history"})
	}
	return c.JSON(updatedHistory)
}

func DeleteHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	err = services.DeleteHistoryByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(fiber.Map{"error": "History not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete history"})
	}
	return c.SendStatus(204)
}
