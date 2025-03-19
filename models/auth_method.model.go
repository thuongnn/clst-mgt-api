package models

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AuthMethod struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	IsActive  bool               `json:"is_active" bson:"is_active" default:"true"`
	Configs   json.RawMessage    `json:"configs,omitempty" bson:"configs,omitempty"`
	CreateAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type AuthMethodResponse struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	IsActive  bool               `json:"is_active" bson:"is_active" default:"true"`
	CreateAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type BasicAuthConfig struct {
	ButtonText string `json:"button_text"`
}

type OAuth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	IssuerURL    string   `json:"issuer_url"`
	RedirectURL  string   `json:"redirect_url"`
	Scopes       []string `json:"scopes"`
	ButtonText   string   `json:"button_text"`
}
