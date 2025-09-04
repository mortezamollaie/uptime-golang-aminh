package controllers

import (
	"strconv"
	"uptime/database"
	"uptime/models"
	"uptime/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateNodeLog(c *fiber.Ctx) error {
	var log models.NodeLog
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	newLog, err := services.CreateNodeLog(&log)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(newLog)
}

func GetAllNodeLogs(c *fiber.Ctx) error {
	logs, err := services.GetAllNodeLogs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}

func GetNodeLog(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	log, err := services.GetNodeLog(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(log)
}

func UpdateNodeLog(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var body models.NodeLog
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	log, err := services.GetNodeLog(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	log.Delay = body.Delay
	log.Status = body.Status
	log.Up = body.Up
	log.Suspended = body.Suspended
	log.Exception = body.Exception

	updatedLog, err := services.UpdateNodeLog(log)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedLog)
}

func DeleteNodeLog(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := services.DeleteNodeLogByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
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
