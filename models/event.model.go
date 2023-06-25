package models

import "k8s.io/apimachinery/pkg/util/json"

type EventType string

const (
	TriggerAll       EventType = "trigger_all"
	TriggerByRuleIds EventType = "trigger_by_rule_ids"
)

type EventMessage struct {
	Type EventType   `json:"event_type"`
	Data interface{} `json:"data"`
}

func (e *EventMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EventMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}
