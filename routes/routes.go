package routes

import (
	"uptime/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	node := api.Group("/nodes")
	node.Get("/with-logs/all", controllers.GetAllNodesWithLogs)
	node.Post("/", controllers.CreateNode)
	node.Get("/", controllers.GetAllNodes)
	node.Get("/:id", controllers.GetNode)
	node.Put("/:id", controllers.UpdateNode)
	node.Delete("/:id", controllers.DeleteNode)

	nodeLogs := api.Group("/node-logs")
	nodeLogs.Post("/", controllers.CreateNodeLog)
	nodeLogs.Get("/", controllers.GetAllNodeLogs)
	nodeLogs.Get("/:id", controllers.GetNodeLog)
	nodeLogs.Put("/:id", controllers.UpdateNodeLog)
	nodeLogs.Delete("/:id", controllers.DeleteNodeLog)

	histories := api.Group("/histories")
	histories.Post("/", controllers.CreateHistory)
	histories.Get("/", controllers.GetAllHistories)
	histories.Get("/:id", controllers.GetHistory)
	histories.Put("/:id", controllers.UpdateHistory)
	histories.Delete("/:id", controllers.DeleteHistory)

	api.Get("/check-uptime", controllers.CheckUptime)
}
