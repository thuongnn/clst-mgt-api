package services

import (
	"github.com/wpcodevo/golang-mongodb/models"
)

type NodeService interface {
	GetRoles() ([]string, error)
	GetRolesByNodeId(string) ([]string, error)
	GetNodesByRoles([]string) ([]*models.DBNode, error)
	GetNodes() ([]*models.DBNode, error)
	CreateNode(*models.DBNode) error
	UpdateByNodeID(string, *models.DBNode) error
	SyncNodes() error
	IsExists(string) (bool, error)
}
