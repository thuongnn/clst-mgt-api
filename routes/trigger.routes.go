package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/controllers"
	"github.com/thuongnn/clst-mgt-api/middleware"
	"github.com/thuongnn/clst-mgt-api/services"
)

type TriggerRouteController struct {
	triggerController controllers.TriggerController
}

func NewTriggerControllerRoute(triggerController controllers.TriggerController) TriggerRouteController {
	return TriggerRouteController{triggerController}
}

func (t *TriggerRouteController) TriggerRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/triggers")
	router.Use(middleware.DeserializeUser(userService))

	router.POST("/all", middleware.AdminOnly(), t.triggerController.TriggerAll)
	router.POST("/", t.triggerController.TriggerByRuleIds)

}
