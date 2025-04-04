package services

import "github.com/thuongnn/clst-mgt-api/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.UserDBResponse, error)
	SignInUser(*models.SignInInput) (*models.UserDBResponse, error)
	SyncOauth2User(*models.SignUpInput) (*models.UserDBResponse, error)
}
