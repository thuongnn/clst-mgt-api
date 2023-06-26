package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/controllers"
	"github.com/thuongnn/clst-mgt-api/middleware"
	"github.com/thuongnn/clst-mgt-api/services"
)

type HistoryScanRouteController struct {
	historyScanController controllers.HistoryScanController
}

func NewHistoryScanControllerRoute(historyScanController controllers.HistoryScanController) HistoryScanRouteController {
	return HistoryScanRouteController{historyScanController}
}

func (r *HistoryScanRouteController) HistoryScanRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/history-scan")
	router.Use(middleware.DeserializeUser(userService))

	router.GET("/:ruleId", r.historyScanController.GetHistoryScanByRuleId)
}
