package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignUpInput struct
type SignUpInput struct {
	Name       string    `json:"name" bson:"name" binding:"required"`
	Username   string    `json:"username" bson:"username" binding:"required"`
	Password   string    `json:"password" bson:"password" binding:"required,min=8"`
	Email      string    `json:"email" bson:"email"`
	Role       string    `json:"role" bson:"role"`
	Verified   bool      `json:"verified" bson:"verified"`
	IsActive   bool      `json:"is_active" bson:"is_active" default:"true"`
	AuthMethod string    `json:"auth_method" bson:"auth_method"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

// SignInInput struct
type SignInInput struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// DBResponse struct
type DBResponse struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Username   string             `json:"username" bson:"username"`
	Password   string             `json:"password" bson:"password"`
	Email      string             `json:"email" bson:"email"`
	Role       string             `json:"role" bson:"role"`
	Verified   bool               `json:"verified" bson:"verified"`
	IsActive   bool               `json:"is_active" bson:"is_active"`
	AuthMethod string             `json:"auth_method" bson:"auth_method"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdateInput struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Role      string    `json:"role,omitempty" bson:"role,omitempty"`
	Verified  bool      `json:"verified,omitempty" bson:"verified,omitempty"`
	IsActive  bool      `json:"is_active,omitempty" bson:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// UserResponse struct
type UserResponse struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Username   string             `json:"username,omitempty" bson:"username,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	Role       string             `json:"role,omitempty" bson:"role,omitempty"`
	Verified   bool               `json:"verified,omitempty" bson:"verified,omitempty"`
	IsActive   bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	AuthMethod string             `json:"auth_method,omitempty" bson:"auth_method,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserClaims struct {
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"email_verified"`
	Name              string   `json:"name"`
	GivenName         string   `json:"given_name"`
	PreferredUsername string   `json:"preferred_username"`
	Nickname          string   `json:"nickname"`
	Groups            []string `json:"groups"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		Role:       user.Role,
		Verified:   user.Verified,
		IsActive:   user.IsActive,
		AuthMethod: user.AuthMethod,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
