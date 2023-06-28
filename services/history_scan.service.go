package services

import (
	"github.com/thuongnn/clst-mgt-api/models"
)

type HistoryScanService interface {
	CreateHistoryScan(historyScan *models.DBHistoryScan) error
	GetHistoryScanByRuleId(ruleId string) ([]*models.DBHistoryScan, error)
	CleanUpHistoryScanByRuleId(ruleId string) error
}
