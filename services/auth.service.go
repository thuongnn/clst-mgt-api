package services

import "github.com/thuongnn/clst-mgt-api/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	SignInUser(*models.SignInInput) (*models.DBResponse, error)
	SyncOauth2User(*models.SignUpInput) (*models.DBResponse, error)
}
