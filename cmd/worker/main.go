package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/thuongnn/clst-mgt-api/config"
	"github.com/thuongnn/clst-mgt-api/handlers"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
	"github.com/thuongnn/clst-mgt-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"k8s.io/client-go/kubernetes"
	"log"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client
	k8sClient   *kubernetes.Clientset

	msgHandler *handlers.MessageHandler
)

func init() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// 👇 Connect to MongoDB
	mongoConnection := options.Client().ApplyURI(appConfig.DBUri)
	mongoClient, err = mongo.Connect(ctx, mongoConnection)
	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("MongoDB successfully connected..")

	// 👇 Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		DB:   0,
		Addr: appConfig.RedisUri,
		//Password: appConfig.RedisPassword,
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	log.Println("Redis successfully connected..")

	// 👇 Connect to K8s cluster
	kubeConfig, err := utils.GetKubeConfig(appConfig.Environment != config.DefaultEnvironment)
	if err != nil {
		panic(err)
	}

	k8sClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}

	if err := utils.K8SHealth(k8sClient, ctx); err != nil {
		panic(err)
	}
	log.Println("Kubernetes API successfully connected..")

	//
	msgHandler = handlers.NewMessageHandler(ctx)
}

func main() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	defer mongoClient.Disconnect(ctx)
	defer redisClient.Close()

	// 👇 Nodes
	nodeCollection := mongoClient.Database(appConfig.DBName).Collection("nodes")
	nodeService := services.NewNodeService(nodeCollection, k8sClient, ctx)

	// 👇 Rules
	ruleCollection := mongoClient.Database(appConfig.DBName).Collection("rules")
	ruleService := services.NewRuleService(ruleCollection, ctx)

	// fw handlers register
	fwHandler := handlers.NewFWHandler(ctx, nodeService, ruleService, redisClient, k8sClient)
	msgHandler.RegisterHandler(models.TriggerAll, fwHandler.HandleScanAllRules)
	msgHandler.RegisterHandler(models.TriggerByRuleIds, fwHandler.HandleScanByRuleIds)

	// starting to handle message from Redis PubSub
	startHandleMessage()
}

func startHandleMessage() {
	// Subscribe to the Topic given
	topic := redisClient.Subscribe(ctx, "rule_triggers")
	channel := topic.Channel()
	for msg := range channel {
		message := &models.EventMessage{}
		// Unmarshal the data into the user
		err := message.UnmarshalBinary([]byte(msg.Payload))
		if err != nil {
			panic(err)
		}

		if errHandler := msgHandler.HandleMessage(message); errHandler != nil {
			log.Println(errHandler)
		}
	}
}
