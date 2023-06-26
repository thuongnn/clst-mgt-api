package services

import (
	"github.com/thuongnn/clst-mgt-api/models"
)

type RuleService interface {
	GetRules(page int, limit int) (*models.RuleListResponse, error)
	GetRulesByRoles(roles []string) ([]*models.DBRule, error)
	GetRulesByIdsAndRoles(ids []string, roles []string) ([]*models.DBRule, error)
	CreateRule(rule *models.DBRule) error
	UpdateRule(id string, rule *models.UpdateRule) error
	DeleteRule(id string) error
}
