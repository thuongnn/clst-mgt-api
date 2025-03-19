package services

import "github.com/thuongnn/clst-mgt-api/models"

type UserService interface {
	FindUsers(params *models.UserSearchParams) (*models.UserListResponse, error)
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByUsername(username string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, data *models.UserUpdate) error
}
