package controllers

import (
	"strconv"
	"uptime/services"

	"github.com/gofiber/fiber/v2"
)

func CreateNode(c *fiber.Ctx) error {
	type Request struct {
		URL string `json:"url"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	node, err := services.CreateNode(body.URL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(node)
}

func GetAllNodes(c *fiber.Ctx) error {
	nodes, err := services.GetAllNodes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(nodes)
}

func GetNode(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	node, err := services.GetNode(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(node)
}

func UpdateNode(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	type Request struct {
		URL string `json:"url"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	node, err := services.UpdateNodeURL(uint(id), body.URL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(node)
}

func DeleteNode(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := services.DeleteNodeByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
