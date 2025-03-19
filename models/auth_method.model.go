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
	ClientID           string   `json:"client_id"`
	ClientSecret       string   `json:"client_secret"`
	WellKnownConfigURL string   `json:"well_known_config_url"`
	IssuerURL          string   `json:"issuer_url"`
	RedirectURL        string   `json:"redirect_url"`
	Scopes             []string `json:"scopes"`
	AdminGroups        []string `json:"admin_groups"`
	ButtonText         string   `json:"button_text"`
}

type WellKnownConfig struct {
	Issuer      string `json:"issuer"`
	AuthURL     string `json:"authorization_endpoint"`
	TokenURL    string `json:"token_endpoint"`
	UserInfoURL string `json:"userinfo_endpoint"`
}
