package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/controllers"
	"github.com/thuongnn/clst-mgt-api/middleware"
	"github.com/thuongnn/clst-mgt-api/services"
)

type SettingRouteController struct {
	authMethodController controllers.AuthMethodController
}

func NewSettingControllerRoute(authMethodController controllers.AuthMethodController) SettingRouteController {
	return SettingRouteController{authMethodController}
}

func (s *SettingRouteController) SettingRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/settings")
	router.Use(middleware.DeserializeUser(userService))
	router.Use(middleware.AdminOnly())

	router.GET("/auth/", s.authMethodController.GetAuthMethods)
	router.GET("/auth/:authMethodId", s.authMethodController.GetAuthMethodById)
	router.POST("/auth/", s.authMethodController.CreateAuthMethod)
	router.PATCH("/auth/:authMethodId", s.authMethodController.UpdateAuthMethod)
	router.DELETE("/auth/:authMethodId", s.authMethodController.DeleteAuthMethod)
}
