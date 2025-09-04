package services

import (
	"errors"
	"uptime/models"
	"uptime/repositories"
)

func CreateNodeLog(log *models.NodeLog) (*models.NodeLog, error) {
	err := repositories.CreateNodeLog(log)
	return log, err
}

func GetAllNodeLogs() ([]models.NodeLog, error) {
	var logs []models.NodeLog
	err := repositories.GetAllNodeLogs(&logs)
	return logs, err
}

func GetNodeLog(id uint) (*models.NodeLog, error) {
	log := &models.NodeLog{}
	err := repositories.GetNodeLogByID(id, log)
	if err != nil {
		return nil, errors.New("NodeLog not found")
	}
	return log, nil
}

func UpdateNodeLog(log *models.NodeLog) (*models.NodeLog, error) {
	err := repositories.UpdateNodeLog(log)
	return log, err
}

func DeleteNodeLogByID(id uint) error {
	log, err := GetNodeLog(id)
	if err != nil {
		return err
	}
	return repositories.DeleteNodeLog(log)
}
