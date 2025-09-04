package services

import (
	"errors"
	"uptime/models"
	"uptime/repositories"
)

func CreateHistory(history *models.History) (*models.History, error) {
	err := repositories.CreateHistory(history)
	return history, err
}

func GetAllHistories() ([]models.History, error) {
	var histories []models.History
	err := repositories.GetAllHistories(&histories)
	return histories, err
}

func GetHistory(id uint) (*models.History, error) {
	history := &models.History{}
	err := repositories.GetHistoryByID(id, history)
	if err != nil {
		return nil, errors.New("History not found")
	}
	return history, nil
}

func UpdateHistory(history *models.History) (*models.History, error) {
	err := repositories.UpdateHistory(history)
	return history, err
}

func DeleteHistoryByID(id uint) error {
	history, err := GetHistory(id)
	if err != nil {
		return err
	}
	return repositories.DeleteHistory(history)
}
