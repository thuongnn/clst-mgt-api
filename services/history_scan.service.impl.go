package services

import (
	"context"
	"github.com/thuongnn/clst-mgt-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type HistoryScanServiceImpl struct {
	historyScanCollection *mongo.Collection
	ctx                   context.Context
}

func (h HistoryScanServiceImpl) CleanUpHistoryScanByRuleId(ruleId string) error {
	obId, _ := primitive.ObjectIDFromHex(ruleId)
	filter := bson.M{"rule_id": obId}

	deleteResult, err := h.historyScanCollection.DeleteMany(h.ctx, filter)
	if err != nil {
		return err
	}

	log.Printf("Deleted %v documents with rule id: %s", deleteResult.DeletedCount, ruleId)
	return nil
}

func (h HistoryScanServiceImpl) CreateHistoryScan(historyScan *models.DBHistoryScan) error {
	filter := bson.M{
		"rule_id":             historyScan.RuleId,
		"node_id":             historyScan.NodeId,
		"destination_address": historyScan.DestinationAddress,
		"destination_port":    historyScan.DestinationPort,
	}

	update := bson.M{"$set": bson.M{
		"node_address":     historyScan.NodeAddress,
		"is_through_proxy": historyScan.IsThroughProxy,
		"node_name":        historyScan.NodeName,
		"error_message":    historyScan.ErrorMessage,
		"status":           historyScan.Status,
		"updated_at":       historyScan.UpdatedAt,
	}}

	upsert := true // create new if not exists
	if _, err := h.historyScanCollection.UpdateOne(h.ctx, filter, update, &options.UpdateOptions{
		Upsert: &upsert,
	}); err != nil {
		return err
	}

	return nil
}

func (h HistoryScanServiceImpl) GetHistoryScanByRuleId(ruleId string) ([]*models.DBHistoryScan, error) {
	obId, _ := primitive.ObjectIDFromHex(ruleId)
	query := bson.M{"rule_id": obId}

	cursor, err := h.historyScanCollection.Find(h.ctx, query)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(h.ctx)

	var records []*models.DBHistoryScan

	for cursor.Next(h.ctx) {
		record := &models.DBHistoryScan{}
		if errDecode := cursor.Decode(record); errDecode != nil {
			return nil, errDecode
		}

		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []*models.DBHistoryScan{}, nil
	}

	return records, nil
}

func NewHistoryScanService(historyScanCollection *mongo.Collection, ctx context.Context) HistoryScanService {
	return &HistoryScanServiceImpl{historyScanCollection, ctx}
}
