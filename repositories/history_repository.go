package repositories

import (
	"uptime/database"
	"uptime/models"
)

func CreateHistory(history *models.History) error {
	return database.DB.Create(history).Error
}

func GetAllHistories(histories *[]models.History) error {
	return database.DB.Find(histories).Error
}

func GetHistoryByID(id uint, history *models.History) error {
	return database.DB.First(history, id).Error
}

func UpdateHistory(history *models.History) error {
	return database.DB.Save(history).Error
}

func DeleteHistory(history *models.History) error {
	return database.DB.Delete(history).Error
}
