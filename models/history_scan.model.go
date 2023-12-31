package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DBHistoryScan struct {
	Id                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RuleId             primitive.ObjectID `json:"rule_id,omitempty" bson:"rule_id,omitempty"`
	NodeId             string             `json:"node_id,omitempty" bson:"node_id,omitempty"`
	NodeName           string             `json:"node_name,omitempty" bson:"node_name,omitempty"`
	NodeAddress        DBNodeAddress      `json:"node_address,omitempty" bson:"node_address,omitempty"`
	DestinationAddress string             `json:"destination_address" bson:"destination_address,omitempty"`
	DestinationPort    string             `json:"destination_port" bson:"destination_port,omitempty"`
	IsThroughProxy     bool               `json:"is_through_proxy,omitempty" bson:"is_through_proxy,omitempty"`
	Status             string             `json:"status,omitempty" bson:"status,omitempty"`
	ErrorMessage       string             `json:"error_message" bson:"error_message,omitempty"`
	UpdatedAt          time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
