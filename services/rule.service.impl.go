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

func (r RuleServiceImpl) IsHistoryScanExits(nodeId, destAddress string, destPort int) (bool, error) {
	filter := bson.M{
		"history_scan": bson.M{
			"$elemMatch": bson.M{
				"node_id":             nodeId,
				"destination_address": destAddress,
				"destination_port":    destPort,
			},
		},
	}

	if err := r.ruleCollection.FindOne(r.ctx, filter).Err(); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
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

func (r RuleServiceImpl) CreateHistoryScan(ruleId primitive.ObjectID, historyScan *models.HistoryScan) error {
	historyScan.UpdatedAt = time.Now()

	isExists, errCheck := r.IsHistoryScanExits(historyScan.NodeId, historyScan.DestinationAddress, historyScan.DestinationPort)
	if errCheck != nil {
		return errCheck
	}

	if !isExists {
		filter := bson.D{{Key: "_id", Value: ruleId}}
		update := bson.M{"$push": bson.M{"history_scan": historyScan}}
		if _, err := r.ruleCollection.UpdateOne(r.ctx, filter, update); err != nil {
			return err
		}
	} else {
		filter := bson.M{
			"_id":                              ruleId,
			"history_scan.node_id":             historyScan.NodeId,
			"history_scan.destination_address": historyScan.DestinationAddress,
			"history_scan.destination_port":    historyScan.DestinationPort,
		}

		update := bson.M{"$set": bson.M{
			"history_scan.$.status":     historyScan.Status,
			"history_scan.$.updated_at": historyScan.UpdatedAt,
		}}

		if _, err := r.ruleCollection.UpdateOne(r.ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

func (r RuleServiceImpl) GetHistoryScanByRuleId(ruleId string) ([]models.HistoryScan, error) {
	obId, _ := primitive.ObjectIDFromHex(ruleId)
	query := bson.M{"_id": obId}

	var rule *models.DBRule
	if err := r.ruleCollection.FindOne(r.ctx, query).Decode(&rule); err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.HistoryScan{}, errors.New("no document with that Id exists")
		}

		return []models.HistoryScan{}, err
	}

	if rule.HistoryScan == nil {
		return []models.HistoryScan{}, nil
	}

	return rule.HistoryScan, nil
}

func NewRuleService(ruleCollection *mongo.Collection, ctx context.Context) RuleService {
	return &RuleServiceImpl{ruleCollection, ctx}
}
