package services

type TriggerService interface {
	TriggerAll() error
	TriggerByRuleIds(ruleIds []string) error
}
