package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBRule struct {
	Id                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status               int                `json:"status,omitempty" bson:"status,omitempty"`
	Roles                []string           `json:"roles,omitempty" bson:"roles,omitempty"`
	Projects             []string           `json:"projects,omitempty" bson:"projects,omitempty"`
	DestinationAddresses []string           `json:"destination_addresses,omitempty" bson:"destination_addresses,omitempty"`
	DestinationPorts     []string           `json:"destination_ports,omitempty" bson:"destination_ports,omitempty"`
	DestinationServices  []string           `json:"destination_services,omitempty" bson:"destination_services,omitempty"`
	IsThroughProxy       bool               `json:"is_through_proxy" bson:"is_through_proxy" default:"false"`
	CR                   []int              `json:"cr,omitempty" bson:"cr,omitempty"`
	IsActive             bool               `json:"is_active" bson:"is_active" default:"true"`
	Description          string             `json:"description,omitempty" bson:"description,omitempty"`
	CreateAt             time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateRule struct {
	Status               int       `json:"status,omitempty" bson:"status,omitempty"`
	Roles                []string  `json:"roles,omitempty" bson:"roles,omitempty"`
	Projects             []string  `json:"projects,omitempty" bson:"projects,omitempty"`
	DestinationAddresses []string  `json:"destination_addresses,omitempty" bson:"destination_addresses,omitempty"`
	DestinationPorts     []string  `json:"destination_ports,omitempty" bson:"destination_ports,omitempty"`
	DestinationServices  []string  `json:"destination_services,omitempty" bson:"destination_services,omitempty"`
	IsThroughProxy       bool      `json:"is_through_proxy,omitempty" bson:"is_through_proxy" default:"false"`
	CR                   []int     `json:"cr,omitempty" bson:"cr,omitempty"`
	IsActive             bool      `json:"is_active,omitempty" bson:"is_active" default:"true"`
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

type RuleSearchParams struct {
	CurrentPage               int    `json:"current_page"`
	PageSize                  int    `json:"page_size"`
	RoleKeyword               string `json:"role_keyword"`
	DestinationAddressKeyword string `json:"destination_address_keyword"`
	CRKeyword                 string `json:"cr_keyword"`
	ProjectKeyword            string `json:"project_keyword"`
}

type Port struct {
	Number   string
	Protocol string
}
