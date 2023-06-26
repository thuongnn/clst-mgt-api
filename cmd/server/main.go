package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/thuongnn/clst-mgt-api/utils"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/config"
	"github.com/thuongnn/clst-mgt-api/controllers"
	"github.com/thuongnn/clst-mgt-api/routes"
	"github.com/thuongnn/clst-mgt-api/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client
	k8sClient   *kubernetes.Clientset

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	// 👇 Create the Nodes Variables
	nodeService         services.NodeService
	NodeController      controllers.NodeController
	nodeCollection      *mongo.Collection
	NodeRouteController routes.NodeRouteController

	// 👇 Create the Rules Variables
	ruleService         services.RuleService
	RuleController      controllers.RuleController
	ruleCollection      *mongo.Collection
	RuleRouteController routes.RuleRouteController

	// 👇 Create the History Scan Variables
	historyScanService         services.HistoryScanService
	HistoryScanController      controllers.HistoryScanController
	historyScanCollection      *mongo.Collection
	HistoryScanRouteController routes.HistoryScanRouteController

	// 👇 Create the Triggers Variables
	triggerService         services.TriggerService
	TriggerController      controllers.TriggerController
	TriggerRouteController routes.TriggerRouteController
)

func init() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoConnection := options.Client().ApplyURI(appConfig.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConnection)
	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	fmt.Println(appConfig.RedisUri)
	fmt.Println(appConfig.RedisPassword)
	redisClient = redis.NewClient(&redis.Options{
		DB:   0,
		Addr: appConfig.RedisUri,
		//Password: appConfig.RedisPassword,
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Println("Redis successfully connected...")

	// Connect to K8s cluster
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
	fmt.Println("Kubernetes API successfully connected...")

	// 👇 Auth
	authCollection = mongoClient.Database(appConfig.DBName).Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService, ctx, authCollection)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	// 👇 Users
	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	// 👇 Nodes
	nodeCollection = mongoClient.Database(appConfig.DBName).Collection("nodes")
	nodeService = services.NewNodeService(nodeCollection, k8sClient, ctx)
	NodeController = controllers.NewNodeController(nodeService)
	NodeRouteController = routes.NewNodeControllerRoute(NodeController)

	// 👇 Rules
	ruleCollection = mongoClient.Database(appConfig.DBName).Collection("rules")
	ruleService = services.NewRuleService(ruleCollection, ctx)
	RuleController = controllers.NewRuleController(ruleService)
	RuleRouteController = routes.NewRuleControllerRoute(RuleController)

	// 👇 History Scan
	historyScanCollection = mongoClient.Database(appConfig.DBName).Collection("history_scan")
	historyScanService = services.NewHistoryScanService(historyScanCollection, ctx)
	HistoryScanController = controllers.NewHistoryScanController(historyScanService)
	HistoryScanRouteController = routes.NewHistoryScanControllerRoute(HistoryScanController)

	// 👇 Triggers
	triggerService = services.NewTriggerService(redisClient, ctx)
	TriggerController = controllers.NewTriggerController(triggerService)
	TriggerRouteController = routes.NewTriggerControllerRoute(TriggerController)

	server = gin.Default()
}

func main() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)
	defer redisClient.Close()

	startGinServer(appConfig)
}

func startGinServer(config config.Config) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	NodeRouteController.NodeRoute(router, userService)
	RuleRouteController.RuleRoute(router, userService)
	HistoryScanRouteController.HistoryScanRoute(router, userService)
	TriggerRouteController.TriggerRoute(router, userService)

	log.Fatal(server.Run(":" + config.Port))
}
