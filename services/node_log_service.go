package services

import (
	"errors"
	"uptime/models"
	"uptime/repositories"
)

func CreateNodeLog(log *models.NodeLog) (*models.NodeLog, error) {
	if log == nil {
		return nil, errors.New("node log cannot be nil")
	}
	
	if log.NodeID == 0 {
		return nil, errors.New("NodeID is required")
	}
	
	err := repositories.CreateNodeLog(log)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func GetAllNodeLogs() ([]models.NodeLog, error) {
	var logs []models.NodeLog
	err := repositories.GetAllNodeLogs(&logs)
	return logs, err
}

func GetNodeLog(id uint) (*models.NodeLog, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	
	log := &models.NodeLog{}
	err := repositories.GetNodeLogByID(id, log)
	if err != nil {
		return nil, errors.New("node log not found")
	}
	return log, nil
}

func UpdateNodeLog(log *models.NodeLog) (*models.NodeLog, error) {
	if log == nil {
		return nil, errors.New("node log cannot be nil")
	}
	
	if log.ID == 0 {
		return nil, errors.New("invalid node log ID")
	}
	
	err := repositories.UpdateNodeLog(log)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func DeleteNodeLogByID(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	
	log, err := GetNodeLog(id)
	if err != nil {
		return err
	}
	return repositories.DeleteNodeLog(log)
}
