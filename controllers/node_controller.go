package controllers

import (
	"net/url"
	"strconv"
	"strings"
	"uptime/services"

	"github.com/gofiber/fiber/v2"
)

// CreateNode creates a new monitoring node
// @Summary Create a new node
// @Description Add a new website/service to monitor
// @Tags nodes
// @Accept json
// @Produce json
// @Param node body object{url=string} true "Node URL"
// @Success 201 {object} map[string]interface{} "Node created successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /nodes [post]
func CreateNode(c *fiber.Ctx) error {
	type Request struct {
		URL string `json:"url"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// Validate URL
	if strings.TrimSpace(body.URL) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "URL cannot be empty"})
	}

	// Parse and validate URL format
	parsedURL, err := url.Parse(body.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid URL format. Must be a valid HTTP/HTTPS URL"})
	}

	node, err := services.CreateNode(body.URL)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return c.Status(409).JSON(fiber.Map{"error": "URL already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create node"})
	}
	return c.Status(201).JSON(node)
}

// GetAllNodes retrieves all monitoring nodes
// @Summary Get all nodes
// @Description Retrieve a list of all monitored websites/services
// @Tags nodes
// @Produce json
// @Success 200 {array} models.Node "List of nodes"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /nodes [get]
func GetAllNodes(c *fiber.Ctx) error {
	nodes, err := services.GetAllNodes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(nodes)
}

func GetNode(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	node, err := services.GetNode(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Node not found"})
	}
	return c.JSON(node)
}

func UpdateNode(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	type Request struct {
		URL string `json:"url"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// Validate URL
	if strings.TrimSpace(body.URL) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "URL cannot be empty"})
	}

	// Parse and validate URL format
	parsedURL, err := url.Parse(body.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid URL format. Must be a valid HTTP/HTTPS URL"})
	}

	node, err := services.UpdateNodeURL(uint(id), body.URL)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(fiber.Map{"error": "Node not found"})
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return c.Status(409).JSON(fiber.Map{"error": "URL already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update node"})
	}
	return c.JSON(node)
}

func DeleteNode(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID parameter is required"})
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	
	err = services.DeleteNodeByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(fiber.Map{"error": "Node not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete node"})
	}
	return c.SendStatus(204)
}
