package services

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/thuongnn/clst-mgt-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func (us *UserServiceImpl) FindUserByUsername(username string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"username": strings.ToLower(username)}
	err := us.userCollection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us UserServiceImpl) FindUsers(params *models.UserSearchParams) (*models.UserListResponse, error) {
	// Set default values for page and limit
	page := params.CurrentPage
	limit := params.PageSize

	filter := bson.M{}

	notEmpty := func(s string) bool {
		return strings.TrimSpace(s) != ""
	}

	if notEmpty(params.NameKeyword) {
		filter["name"] = bson.M{"$regex": params.NameKeyword, "$options": "i"}
	}

	if notEmpty(params.UsernameKeyword) {
		filter["username"] = bson.M{"$regex": params.UsernameKeyword, "$options": "i"}
	}

	if notEmpty(params.EmailKeyword) {
		filter["email"] = bson.M{"$regex": params.EmailKeyword, "$options": "i"}
	}

	// 👇 Calculate the total number of pages
	count, err := us.userCollection.CountDocuments(us.ctx, filter)
	if err != nil {
		return nil, err
	}

	// In case there are no documents matching the filter
	if count == 0 {
		return &models.UserListResponse{
			Data:       []*models.UserResponse{},
			Pagination: &models.Pagination{},
		}, nil
	}

	// Initialize totalPages variable with count & limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	// Check page is not greater than totalPages, get minimum value between page and totalPages
	if page > totalPages {
		page = totalPages
	}

	// build find options with:
	// + Limit	: limit number of documents
	// + Skip	: skip documents in rage for pagination
	// + Sort	: sort by updated_at
	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64((page - 1) * limit))
	opt.SetSort(bson.M{"username": 1})

	// 👇 Starting to finding from DB with filter & pagination
	cursor, err := us.userCollection.Find(us.ctx, filter, &opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(us.ctx)

	// Parsing result
	var users []*models.UserResponse
	for cursor.Next(us.ctx) {
		user := &models.UserResponse{}
		if errDecode := cursor.Decode(user); errDecode != nil {
			return nil, errDecode
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &models.UserListResponse{
		Data: users,
		Pagination: &models.Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			PageSize:    limit,
			TotalCount:  int(count),
		},
	}, nil
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := us.userCollection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.userCollection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) UpdateUserById(id string, data *models.UserUpdate) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	updateQuery := bson.D{{Key: "_id", Value: obId}}
	updateData := bson.D{{"$set", bson.D{{"is_active", data.IsActive}, {"updated_at", time.Now()}}}}
	res := us.userCollection.FindOneAndUpdate(us.ctx, updateQuery, updateData)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
