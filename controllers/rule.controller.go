package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
	"net/http"
	"strconv"
	"strings"
)

type RuleController struct {
	ruleService services.RuleService
}

func NewRuleController(ruleService services.RuleService) RuleController {
	return RuleController{ruleService}
}

func (rc *RuleController) CreateRule(ctx *gin.Context) {
	var rule *models.DBRule

	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := rc.ruleService.CreateRule(rule); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (rc *RuleController) UpdateRule(ctx *gin.Context) {
	ruleId := ctx.Param("ruleId")

	var rule *models.UpdateRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := rc.ruleService.UpdateRule(ruleId, rule); err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (rc *RuleController) GetRules(ctx *gin.Context) {
	var currentPage = ctx.DefaultQuery("current_page", "1")
	var pageSize = ctx.DefaultQuery("page_size", "10")

	intCurrentPage, err := strconv.Atoi(currentPage)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	result, err := rc.ruleService.GetRules(&models.RuleSearchParams{
		CurrentPage:               intCurrentPage,
		PageSize:                  intPageSize,
		RoleKeyword:               ctx.Query("role_keyword"),
		DestinationAddressKeyword: ctx.Query("destination_address_keyword"),
		CRKeyword:                 ctx.Query("cr_keyword"),
		ProjectKeyword:            ctx.Query("project_keyword"),
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"data":       result.Data,
		"total":      result.Pagination.TotalCount,
		"pagination": result.Pagination,
	})
}

func (rc *RuleController) DeleteRule(ctx *gin.Context) {
	ruleId := ctx.Param("ruleId")

	if err := rc.ruleService.DeleteRule(ruleId); err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
