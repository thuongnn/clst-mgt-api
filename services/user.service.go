package services

import "github.com/thuongnn/clst-mgt-api/models"

type UserService interface {
	FindUsers(params *models.UserSearchParams) (*models.UserListResponse, error)
	FindUserById(id string) (*models.UserDBResponse, error)
	FindUserByUsername(username string) (*models.UserDBResponse, error)
	FindUserByEmail(email string) (*models.UserDBResponse, error)
	UpdateUserById(id string, data *models.UserUpdate) error
}
