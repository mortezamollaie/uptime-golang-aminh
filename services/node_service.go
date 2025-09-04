package services

import (
	"errors"
	"uptime/models"
	"uptime/repositories"
)

func CreateNode(url string) (*models.Node, error) {
	node := &models.Node{URL: url}
	err := repositories.CreateNode(node)
	return node, err
}

func GetAllNodes() ([]models.Node, error) {
	var nodes []models.Node
	err := repositories.GetAllNodes(&nodes)
	return nodes, err
}

func GetNode(id uint) (*models.Node, error) {
	node := &models.Node{}
	err := repositories.GetNodeByID(id, node)
	if err != nil {
		return nil, errors.New("Node not found")
	}
	return node, nil
}

func UpdateNodeURL(id uint, newURL string) (*models.Node, error) {
	node, err := GetNode(id)
	if err != nil {
		return nil, err
	}
	node.URL = newURL
	err = repositories.UpdateNode(node)
	return node, err
}

func DeleteNodeByID(id uint) error {
	node, err := GetNode(id)
	if err != nil {
		return err
	}
	return repositories.DeleteNode(node)
}
