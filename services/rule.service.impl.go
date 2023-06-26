package services

import (
	"context"
	"errors"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

type RuleServiceImpl struct {
	ruleCollection *mongo.Collection
	ctx            context.Context
}

func (r RuleServiceImpl) GetRules(page int, limit int) (*models.RuleListResponse, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	// Initialize totalPages variable
	totalPages := 0

	// Calculate the total number of pages
	count, err := r.ruleCollection.CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if count > 0 {
		totalPages = int(math.Ceil(float64(count) / float64(limit)))
	}

	if page > totalPages {
		page = totalPages
	}

	skip := (page - 1) * limit

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))
	opt.SetSort(bson.M{"created_at": -1})

	cursor, err := r.ruleCollection.Find(r.ctx, bson.M{}, &opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.ctx)

	var rules []*models.DBRule
	for cursor.Next(r.ctx) {
		rule := &models.DBRule{}
		if errDecode := cursor.Decode(rule); errDecode != nil {
			return nil, errDecode
		}

		rules = append(rules, rule)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &models.RuleListResponse{
		Data: rules,
		Pagination: &models.Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			PageSize:    limit,
			TotalCount:  int(count),
		},
	}, nil
}

func (r RuleServiceImpl) GetRulesByRoles(roles []string) ([]*models.DBRule, error) {
	filter := bson.M{
		"roles": bson.M{"$in": roles},
	}

	var rules []*models.DBRule
	cursor, err := r.ruleCollection.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.ctx)

	for cursor.Next(r.ctx) {
		var rule = &models.DBRule{}
		if errDecode := cursor.Decode(rule); errDecode != nil {
			return nil, errDecode
		}
		rules = append(rules, rule)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (r RuleServiceImpl) GetRulesByIdsAndRoles(ids []string, roles []string) ([]*models.DBRule, error) {
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs[i] = objectID
	}

	filter := bson.M{
		"_id":   bson.M{"$in": objectIDs},
		"roles": bson.M{"$in": roles},
	}

	var rules []*models.DBRule
	cursor, err := r.ruleCollection.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.ctx)

	for cursor.Next(r.ctx) {
		var rule = &models.DBRule{}
		if errDecode := cursor.Decode(rule); errDecode != nil {
			return nil, errDecode
		}
		rules = append(rules, rule)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (r RuleServiceImpl) CreateRule(rule *models.DBRule) error {
	rule.CreateAt = time.Now()
	rule.UpdatedAt = rule.CreateAt
	rule.Status = utils.StatusUnknown
	rule.IsActive = true

	_, err := r.ruleCollection.InsertOne(r.ctx, rule)
	return err
}

func (r RuleServiceImpl) UpdateRule(id string, rule *models.UpdateRule) error {
	rule.UpdatedAt = time.Now()

	doc, err := utils.ToDoc(rule)
	if err != nil {
		return err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	updateQuery := bson.D{{Key: "_id", Value: obId}}
	updateData := bson.D{{Key: "$set", Value: doc}}
	res := r.ruleCollection.FindOneAndUpdate(r.ctx, updateQuery, updateData)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (r RuleServiceImpl) DeleteRule(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := r.ruleCollection.DeleteOne(r.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Ids exists")
	}

	return nil
}

func NewRuleService(ruleCollection *mongo.Collection, ctx context.Context) RuleService {
	return &RuleServiceImpl{ruleCollection, ctx}
}
