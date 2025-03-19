package services

import (
	"context"
	"errors"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AuthMethodServiceImpl struct {
	authMethodCollection *mongo.Collection
	ctx                  context.Context
}

func (a AuthMethodServiceImpl) GetAuthMethods() ([]*models.AuthMethod, error) {
	filter := bson.M{}

	cursor, err := a.authMethodCollection.Find(a.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(a.ctx)

	var authMethods []*models.AuthMethod
	for cursor.Next(a.ctx) {
		var authMethod models.AuthMethod
		if err := cursor.Decode(&authMethod); err != nil {
			return nil, err
		}
		authMethods = append(authMethods, &authMethod)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return authMethods, nil
}

func (a AuthMethodServiceImpl) GetAuthMethodById(id string) (*models.AuthMethod, error) {
	obId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": obId}

	res := a.authMethodCollection.FindOne(a.ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var authMethod = &models.AuthMethod{}
	if errDecode := res.Decode(authMethod); errDecode != nil {
		return nil, errDecode
	}

	return authMethod, nil
}

func (a AuthMethodServiceImpl) GetActiveAuthMethods() ([]*models.AuthMethod, error) {
	filter := bson.M{"is_active": true}

	cursor, err := a.authMethodCollection.Find(a.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(a.ctx)

	var authMethods []*models.AuthMethod
	for cursor.Next(a.ctx) {
		var authMethod models.AuthMethod
		if err := cursor.Decode(&authMethod); err != nil {
			return nil, err
		}
		authMethods = append(authMethods, &authMethod)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return authMethods, nil
}

func (a AuthMethodServiceImpl) CountAuthMethods() (int64, error) {
	count, err := a.authMethodCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (a AuthMethodServiceImpl) CreateAuthMethod(authMethod *models.AuthMethod) error {
	authMethod.CreateAt = time.Now()
	authMethod.UpdatedAt = authMethod.CreateAt
	authMethod.IsActive = true

	_, err := a.authMethodCollection.InsertOne(a.ctx, authMethod)
	return err
}

func (a AuthMethodServiceImpl) UpdateAuthMethod(id string, authMethod *models.AuthMethod) error {
	authMethod.UpdatedAt = time.Now()

	doc, err := utils.ToDoc(authMethod)
	if err != nil {
		return err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	updateQuery := bson.D{{Key: "_id", Value: obId}}
	updateData := bson.D{{Key: "$set", Value: doc}}
	res := a.authMethodCollection.FindOneAndUpdate(a.ctx, updateQuery, updateData)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (a AuthMethodServiceImpl) DeleteAuthMethod(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := a.authMethodCollection.DeleteOne(a.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Ids exists")
	}

	return nil
}

func NewAuthMethodService(authMethodCollection *mongo.Collection, ctx context.Context) AuthMethodService {
	return &AuthMethodServiceImpl{authMethodCollection, ctx}
}
