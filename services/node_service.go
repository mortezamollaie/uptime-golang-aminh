package services

import (
	"errors"
	"net/url"
	"strings"
	"uptime/models"
	"uptime/repositories"
)

func CreateNode(nodeURL string) (*models.Node, error) {
	// Validate URL
	if strings.TrimSpace(nodeURL) == "" {
		return nil, errors.New("URL cannot be empty")
	}

	// Parse and validate URL format
	parsedURL, err := url.Parse(nodeURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, errors.New("invalid URL format, must be a valid HTTP/HTTPS URL")
	}

	node := &models.Node{URL: nodeURL}
	err = repositories.CreateNode(node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func GetAllNodes() ([]models.Node, error) {
	var nodes []models.Node
	err := repositories.GetAllNodes(&nodes)
	return nodes, err
}

func GetNode(id uint) (*models.Node, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	
	node := &models.Node{}
	err := repositories.GetNodeByID(id, node)
	if err != nil {
		return nil, errors.New("node not found")
	}
	return node, nil
}

func UpdateNodeURL(id uint, newURL string) (*models.Node, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	
	// Validate URL
	if strings.TrimSpace(newURL) == "" {
		return nil, errors.New("URL cannot be empty")
	}

	// Parse and validate URL format
	parsedURL, err := url.Parse(newURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, errors.New("invalid URL format, must be a valid HTTP/HTTPS URL")
	}

	node, err := GetNode(id)
	if err != nil {
		return nil, err
	}
	
	node.URL = newURL
	err = repositories.UpdateNode(node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func DeleteNodeByID(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	
	node, err := GetNode(id)
	if err != nil {
		return err
	}
	return repositories.DeleteNode(node)
}
