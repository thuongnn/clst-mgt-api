package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/thuongnn/clst-mgt-api/models"
)

type TriggerServiceImpl struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (t TriggerServiceImpl) TriggerAll() error {
	// Publish a generated user to the rule_triggers channel
	return t.redisClient.Publish(t.ctx, "rule_triggers", &models.EventMessage{
		Type: models.TriggerAll,
		Data: []string{},
	}).Err()
}

func (t TriggerServiceImpl) TriggerByRuleIds(ruleIds []string) error {
	// Publish a generated user to the rule_triggers channel
	return t.redisClient.Publish(t.ctx, "rule_triggers", &models.EventMessage{
		Type: models.TriggerByRuleIds,
		Data: ruleIds,
	}).Err()
}

func NewTriggerService(redisClient *redis.Client, ctx context.Context) TriggerService {
	return &TriggerServiceImpl{redisClient, ctx}
}
