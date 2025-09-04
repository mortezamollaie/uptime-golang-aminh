package repositories

import (
	"uptime/database"
	"uptime/models"
)

func CreateNode(node *models.Node) error {
	return database.DB.Create(node).Error
}

func GetAllNodes(nodes *[]models.Node) error {
	return database.DB.Find(nodes).Error
}

func GetNodeByID(id uint, node *models.Node) error {
	return database.DB.First(node, id).Error
}

func UpdateNode(node *models.Node) error {
	return database.DB.Save(node).Error
}

func DeleteNode(node *models.Node) error {
	return database.DB.Delete(node).Error
}
