package services

import (
	"errors"
	"uptime/models"
	"uptime/repositories"
)

func CreateHistory(history *models.History) (*models.History, error) {
	if history == nil {
		return nil, errors.New("history cannot be nil")
	}
	
	if history.NodeID == 0 {
		return nil, errors.New("NodeID is required")
	}
	
	err := repositories.CreateHistory(history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func GetAllHistories() ([]models.History, error) {
	var histories []models.History
	err := repositories.GetAllHistories(&histories)
	return histories, err
}

func GetHistory(id uint) (*models.History, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	
	history := &models.History{}
	err := repositories.GetHistoryByID(id, history)
	if err != nil {
		return nil, errors.New("history not found")
	}
	return history, nil
}

func UpdateHistory(history *models.History) (*models.History, error) {
	if history == nil {
		return nil, errors.New("history cannot be nil")
	}
	
	if history.ID == 0 {
		return nil, errors.New("invalid history ID")
	}
	
	err := repositories.UpdateHistory(history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func DeleteHistoryByID(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	
	history, err := GetHistory(id)
	if err != nil {
		return err
	}
	return repositories.DeleteHistory(history)
}
