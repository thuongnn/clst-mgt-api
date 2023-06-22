package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-mongodb/controllers"
	"github.com/wpcodevo/golang-mongodb/middleware"
	"github.com/wpcodevo/golang-mongodb/services"
)

type NodeRouteController struct {
	nodeController controllers.NodeController
}

func NewNodeControllerRoute(nodeController controllers.NodeController) NodeRouteController {
	return NodeRouteController{nodeController}
}

func (r *NodeRouteController) NodeRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/nodes")
	router.Use(middleware.DeserializeUser(userService))

	router.GET("/", r.nodeController.GetNodes)
	router.GET("/sync", r.nodeController.SyncNodes)
	router.GET("/roles", r.nodeController.GetRoles)
	router.GET("/roles/:nodeId", r.nodeController.GetRolesByNodeId)
	//router.PATCH("/:postId", r.postController.UpdatePost)
	//router.DELETE("/:postId", r.postController.DeletePost)
}
