package services

import "github.com/wpcodevo/golang-mongodb/models"

type RuleService interface {
	GetRules(page int, limit int) (*models.RuleListResponse, error)
	GetRulesByIdsAndRoles(ids []string, roles []string) ([]*models.DBRule, error)
	CreateRule(rule *models.DBRule) error
	UpdateRule(id string, rule *models.UpdateRule) error
	DeleteRule(id string) error
	CreateHistoryScan(ruleId string, historyScan *models.HistoryScan) error
	GetHistoryScanByRuleId(ruleId string) ([]models.HistoryScan, error)
	IsHistoryScanExitsWithNodeId(nodeId string) (bool, error)
}
