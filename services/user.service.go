package services

import "github.com/thuongnn/clst-mgt-api/models"

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByUsername(username string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error)
}
