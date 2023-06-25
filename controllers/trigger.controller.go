package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/services"
	"net/http"
)

type TriggerController struct {
	triggerService services.TriggerService
}

func NewTriggerController(triggerService services.TriggerService) TriggerController {
	return TriggerController{triggerService}
}

func (tc *TriggerController) TriggerAll(ctx *gin.Context) {
	if err := tc.triggerService.TriggerAll(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
}

func (tc *TriggerController) TriggerByRuleIds(ctx *gin.Context) {
	type ParseData struct {
		RuleIds []string `json:"rule_ids"`
	}

	var parseData *ParseData
	if err := ctx.ShouldBindJSON(&parseData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := tc.triggerService.TriggerByRuleIds(parseData.RuleIds); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
}
