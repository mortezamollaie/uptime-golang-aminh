package controllers

import (
	"strconv"
	"uptime/models"
	"uptime/services"

	"github.com/gofiber/fiber/v2"
)

func CreateHistory(c *fiber.Ctx) error {
	var history models.History
	if err := c.BodyParser(&history); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	newHistory, err := services.CreateHistory(&history)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(newHistory)
}

func GetAllHistories(c *fiber.Ctx) error {
	histories, err := services.GetAllHistories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(histories)
}

func GetHistory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	history, err := services.GetHistory(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(history)
}

func UpdateHistory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var body models.History
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	history, err := services.GetHistory(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	// Update fields
	history.Delay = body.Delay
	history.Status = body.Status
	history.Up = body.Up
	history.Suspended = body.Suspended
	history.Exception = body.Exception

	updatedHistory, err := services.UpdateHistory(history)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedHistory)
}

func DeleteHistory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := services.DeleteHistoryByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
