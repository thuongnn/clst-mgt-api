package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type NodeServiceImpl struct {
	nodeCollection *mongo.Collection
	k8sClient      *kubernetes.Clientset
	ctx            context.Context
}

func (n NodeServiceImpl) GetRolesByNodeId(nodeId string) ([]string, error) {
	query := bson.M{"node_id": nodeId}

	var node *models.DBNode
	if err := n.nodeCollection.FindOne(n.ctx, query).Decode(&node); err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, errors.New("no document with that Id exists")
		}

		return []string{}, err
	}

	return node.Roles, nil
}

func (n NodeServiceImpl) GetRoles() ([]string, error) {
	pipeline := []bson.M{
		{
			"$unwind": "$roles",
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"roles": bson.M{"$addToSet": "$roles"},
			},
		},
		{
			"$project": bson.M{
				"_id":   0,
				"roles": 1,
			},
		},
	}

	cursor, err := n.nodeCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(n.ctx)

	var result bson.M
	if cursor.Next(n.ctx) {
		errDecode := cursor.Decode(&result)
		if errDecode != nil {
			return []string{}, errDecode
		}
	}

	rolesRaw, ok := result["roles"].(primitive.A)
	if !ok {
		return []string{}, fmt.Errorf("Conversion to []interface{} failed. ")
	}

	var roles []string
	for _, role := range rolesRaw {
		if roleString, ok := role.(string); ok {
			roles = append(roles, roleString)
		} else {
			return []string{}, fmt.Errorf("Conversion to string failed for role: ")
		}
	}

	return roles, nil
}

func (n NodeServiceImpl) IsExists(nodeId string) (bool, error) {
	var dbNode *models.DBNode
	query := bson.M{"node_id": nodeId}
	if err := n.nodeCollection.FindOne(n.ctx, query).Decode(&dbNode); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (n NodeServiceImpl) CreateNode(node *models.DBNode) error {
	node.CreateAt = time.Now()
	node.UpdatedAt = node.CreateAt
	if _, err := n.nodeCollection.InsertOne(n.ctx, &node); err != nil {
		return err
	}
	return nil
}

func (n NodeServiceImpl) UpdateByNodeID(nodeId string, node *models.DBNode) error {
	node.UpdatedAt = time.Now()

	doc, err := utils.ToDoc(node)
	if err != nil {
		return err
	}

	updateQuery := bson.D{{Key: "node_id", Value: nodeId}}
	updateData := bson.D{{Key: "$set", Value: doc}}
	res := n.nodeCollection.FindOneAndUpdate(n.ctx, updateQuery, updateData)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (n NodeServiceImpl) GetNodes() ([]*models.DBNode, error) {
	query := bson.M{}
	cursor, err := n.nodeCollection.Find(n.ctx, query)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(n.ctx)

	var nodes []*models.DBNode

	for cursor.Next(n.ctx) {
		node := &models.DBNode{}
		if errDecode := cursor.Decode(node); errDecode != nil {
			return nil, errDecode
		}

		nodes = append(nodes, node)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []*models.DBNode{}, nil
	}

	return nodes, nil
}

func (n NodeServiceImpl) GetNodesByRoles(labels []string) ([]*models.DBNode, error) {
	query := bson.M{"roles": bson.M{"$in": labels}}
	cursor, err := n.nodeCollection.Find(n.ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(n.ctx)

	var nodes []*models.DBNode
	for cursor.Next(n.ctx) {
		node := &models.DBNode{}
		if errDecode := cursor.Decode(node); errDecode != nil {
			return nil, errDecode
		}

		nodes = append(nodes, node)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []*models.DBNode{}, nil
	}

	return nodes, nil
}

func (n NodeServiceImpl) SyncNodes() error {
	nodes, err := n.k8sClient.CoreV1().Nodes().List(n.ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Start scanning...
	for _, node := range nodes.Items {
		var data models.K8sNode

		data.NodeId = string(node.UID)
		data.Name = node.Name

		var roles []string
		for label, _ := range node.Labels {
			if strings.Contains(label, "node-role.kubernetes.io") {
				roles = append(roles, strings.ReplaceAll(label, "node-role.kubernetes.io/", ""))
			}
		}
		data.Roles = roles

		for _, v := range node.Status.Addresses {
			if v.Type == "ExternalIP" {
				data.Address.ExternalIP = v.Address
			}
			if v.Type == "InternalIP" {
				data.Address.InternalIP = v.Address
			}
			if v.Type == "Hostname" {
				data.Address.Hostname = v.Address
			}
		}

		isNodeExists, errCheck := n.IsExists(data.NodeId)
		if errCheck != nil {
			return errCheck
		}

		// Create new node if not exists in mongodb
		if !isNodeExists {
			if errCreate := n.CreateNode(models.ToDBNode(data)); errCreate != nil {
				return errCreate
			}
			continue
		}

		// Update current node in mongodb if fetch from k8s
		if isNodeExists {
			if errUpdate := n.UpdateByNodeID(data.NodeId, models.ToDBNode(data)); errUpdate != nil {
				return errUpdate
			}
		}
	}

	return nil
}

func NewNodeService(nodeCollection *mongo.Collection, k8sClient *kubernetes.Clientset, ctx context.Context) NodeService {
	return &NodeServiceImpl{nodeCollection, k8sClient, ctx}
}
