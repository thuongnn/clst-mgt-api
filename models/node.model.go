package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DBNode struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NodeId      string             `json:"node_id,omitempty" bson:"node_id,omitempty" binding:"required"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	Roles       []string           `json:"roles,omitempty" bson:"roles,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Address     DBNodeAddress      `json:"address,omitempty" bson:"address,omitempty"`
	CreateAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DBNodeAddress struct {
	InternalIP string `json:"internal_ip,omitempty" bson:"internal_ip,omitempty"`
	ExternalIP string `json:"external_ip,omitempty" bson:"external_ip,omitempty"`
	Hostname   string `json:"hostname,omitempty" bson:"hostname,omitempty"`
}

type K8sNode struct {
	NodeId    string        `json:"node_id,omitempty" bson:"node_id,omitempty" binding:"required"`
	Name      string        `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	Roles     []string      `json:"roles,omitempty" bson:"roles,omitempty"`
	Address   DBNodeAddress `json:"address,omitempty" bson:"address,omitempty"`
	CreateAt  time.Time     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func ToDBNode(node K8sNode) *DBNode {
	return &DBNode{
		NodeId:    node.NodeId,
		Name:      node.Name,
		Roles:     node.Roles,
		Address:   node.Address,
		CreateAt:  node.CreateAt,
		UpdatedAt: node.UpdatedAt,
	}
}
