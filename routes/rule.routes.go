package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/controllers"
	"github.com/thuongnn/clst-mgt-api/middleware"
	"github.com/thuongnn/clst-mgt-api/services"
)

type RuleRouteController struct {
	ruleController controllers.RuleController
}

func NewRuleControllerRoute(ruleController controllers.RuleController) RuleRouteController {
	return RuleRouteController{ruleController}
}

func (r *RuleRouteController) RuleRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/rules")
	router.Use(middleware.DeserializeUser(userService))

	router.GET("/", r.ruleController.GetRules)
	router.POST("/", r.ruleController.CreateRule)
	router.PATCH("/:ruleId", r.ruleController.UpdateRule)
	router.DELETE("/:ruleId", r.ruleController.DeleteRule)
	router.GET("/:ruleId/history", r.ruleController.GetHistoryScanByRuleId)
}
