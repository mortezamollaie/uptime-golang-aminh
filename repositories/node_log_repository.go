package repositories

import (
	"uptime/database"
	"uptime/models"
)

func CreateNodeLog(log *models.NodeLog) error {
	return database.DB.Create(log).Error
}

func GetAllNodeLogs(logs *[]models.NodeLog) error {
	return database.DB.Find(logs).Error
}

func GetNodeLogByID(id uint, log *models.NodeLog) error {
	return database.DB.First(log, id).Error
}

func UpdateNodeLog(log *models.NodeLog) error {
	return database.DB.Save(log).Error
}

func DeleteNodeLog(log *models.NodeLog) error {
	return database.DB.Delete(log).Error
}
