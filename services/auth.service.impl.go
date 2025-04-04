package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthService(collection *mongo.Collection, ctx context.Context) AuthService {
	return &AuthServiceImpl{collection, ctx}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.UserDBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	user.Verified = false
	user.Role = "anonymous"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	res, err := uc.collection.InsertOne(uc.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that username already exist")
		}
		return nil, err
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"username": 1}, Options: opt}

	if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, errors.New("could not create index for username")
	}

	var newUser *models.UserDBResponse
	query := bson.M{"_id": res.InsertedID}

	err = uc.collection.FindOne(uc.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.UserDBResponse, error) {
	return nil, nil
}

func (uc *AuthServiceImpl) SyncOauth2User(user *models.SignUpInput) (*models.UserDBResponse, error) {
	user.UpdatedAt = time.Now()
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	filter := bson.M{"$or": []bson.M{{"email": user.Email}, {"username": user.Username}}}

	var existingUser *models.UserDBResponse
	err := uc.collection.FindOne(uc.ctx, filter).Decode(&existingUser)

	if err != nil {
		// If user does not exist → Create new
		if err == mongo.ErrNoDocuments {
			user.IsActive = true
			user.CreatedAt = user.UpdatedAt
			res, err := uc.collection.InsertOne(uc.ctx, user)
			if err != nil {
				if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
					return nil, errors.New("user with that username or email already exists")
				}
				return nil, err
			}

			var newUser *models.UserDBResponse
			query := bson.M{"_id": res.InsertedID}

			err = uc.collection.FindOne(uc.ctx, query).Decode(&newUser)
			if err != nil {
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}

	// If user already exists, update information
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"verified":   user.Verified,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = uc.collection.UpdateOne(uc.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	err = uc.collection.FindOne(uc.ctx, filter).Decode(&existingUser)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}
