package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/services"
	"net/http"
	"strings"
)

type HistoryScanController struct {
	historyScanService services.HistoryScanService
}

func NewHistoryScanController(historyScanService services.HistoryScanService) HistoryScanController {
	return HistoryScanController{historyScanService}
}

func (hsc *HistoryScanController) GetHistoryScanByRuleId(ctx *gin.Context) {
	ruleId := ctx.Param("ruleId")

	historyScan, err := hsc.historyScanService.GetHistoryScanByRuleId(ruleId)
	if err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(historyScan), "data": historyScan})
}
