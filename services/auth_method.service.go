package services

import (
	"github.com/thuongnn/clst-mgt-api/models"
)

type AuthMethodService interface {
	GetAuthMethods() ([]*models.AuthMethod, error)
	GetAuthMethodById(id string) (*models.AuthMethod, error)
	GetActiveAuthMethods() ([]*models.AuthMethod, error)
	CountAuthMethods() (int64, error)
	CreateAuthMethod(authMethod *models.AuthMethod) error
	UpdateAuthMethod(id string, authMethod *models.AuthMethod) error
	DeleteAuthMethod(id string) error
}
