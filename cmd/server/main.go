package main

import (
	"context"
	"fmt"
	"github.com/wpcodevo/golang-mongodb/utils"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-mongodb/config"
	"github.com/wpcodevo/golang-mongodb/controllers"
	"github.com/wpcodevo/golang-mongodb/routes"
	"github.com/wpcodevo/golang-mongodb/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	k8sClient   *kubernetes.Clientset

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	// 👇 Create the Post Variables
	postService         services.PostService
	PostController      controllers.PostController
	postCollection      *mongo.Collection
	PostRouteController routes.PostRouteController

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

	// Collections
	authCollection = mongoClient.Database(appConfig.DBName).Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService, ctx, authCollection)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	// 👇 Instantiate the Constructors
	postCollection = mongoClient.Database(appConfig.DBName).Collection("posts")
	postService = services.NewPostService(postCollection, ctx)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewPostControllerRoute(PostController)

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

	server = gin.Default()
}

func main() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	// startGinServer(config)
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

	// 👇 Post Route
	PostRouteController.PostRoute(router)

	log.Fatal(server.Run(":" + config.Port))
}
