package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DBRule struct {
	Id                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status               int                `json:"status,omitempty" bson:"status,omitempty"`
	Roles                []string           `json:"roles,omitempty" bson:"roles,omitempty"`
	Projects             []string           `json:"projects,omitempty" bson:"projects,omitempty"`
	DestinationAddresses []string           `json:"destination_addresses,omitempty" bson:"destination_addresses,omitempty"`
	DestinationPorts     []int              `json:"destination_ports,omitempty" bson:"destination_ports,omitempty"`
	DestinationServices  []string           `json:"destination_services,omitempty" bson:"destination_services,omitempty"`
	CR                   []int              `json:"cr,omitempty" bson:"cr,omitempty"`
	IsActive             bool               `json:"is_active,omitempty" bson:"is_active,omitempty" default:"true"`
	Description          string             `json:"description,omitempty" bson:"description,omitempty"`
	HistoryScan          []HistoryScan      `json:"history_scan,omitempty" bson:"history_scan,omitempty"`
	CreateAt             time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateRule struct {
	Status               int       `json:"status,omitempty" bson:"status,omitempty"`
	Roles                []string  `json:"roles,omitempty" bson:"roles,omitempty"`
	Projects             []string  `json:"projects,omitempty" bson:"projects,omitempty"`
	DestinationAddresses []string  `json:"destination_addresses,omitempty" bson:"destination_addresses,omitempty"`
	DestinationPorts     []int     `json:"destination_ports,omitempty" bson:"destination_ports,omitempty"`
	DestinationServices  []string  `json:"destination_services,omitempty" bson:"destination_services,omitempty"`
	CR                   []int     `json:"cr,omitempty" bson:"cr,omitempty"`
	IsActive             bool      `json:"is_active,omitempty" bson:"is_active,omitempty" default:"true"`
	Description          string    `json:"description,omitempty" bson:"description,omitempty"`
	UpdatedAt            time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	PageSize    int `json:"page_size"`
	TotalCount  int `json:"total_count"`
}

type RuleListResponse struct {
	Data       []*DBRule   `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

type DeleteRule struct {
	ids []string
}

type HistoryScan struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	NodeId    string    `json:"node_id,omitempty" bson:"node_id,omitempty"`
	Status    string    `json:"status,omitempty" bson:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
