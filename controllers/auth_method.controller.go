package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
	"net/http"
	"strings"
)

type AuthMethodController struct {
	authMethodService services.AuthMethodService
}

func NewAuthMethodController(authMethodService services.AuthMethodService) AuthMethodController {
	return AuthMethodController{authMethodService}
}

func (sc *AuthMethodController) CreateAuthMethod(ctx *gin.Context) {
	var authMethod *models.AuthMethod

	if err := ctx.ShouldBindJSON(&authMethod); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Check the number of AuthMethods in the database
	count, err := sc.authMethodService.CountAuthMethods()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Maximum limit of 3 AuthMethods
	if count >= 3 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Only create up to 3 authentication methods"})
		return
	}

	if err := sc.authMethodService.CreateAuthMethod(authMethod); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (sc *AuthMethodController) UpdateAuthMethod(ctx *gin.Context) {
	authMethodId := ctx.Param("authMethodId")

	var authMethod *models.AuthMethod
	if err := ctx.ShouldBindJSON(&authMethod); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := sc.authMethodService.UpdateAuthMethod(authMethodId, authMethod); err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (sc *AuthMethodController) GetAuthMethods(ctx *gin.Context) {
	result, err := sc.authMethodService.GetAuthMethods()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

func (sc *AuthMethodController) GetAuthMethodById(ctx *gin.Context) {
	authMethodId := ctx.Param("authMethodId")

	result, err := sc.authMethodService.GetAuthMethodById(authMethodId)
	if err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

func (sc *AuthMethodController) DeleteAuthMethod(ctx *gin.Context) {
	authMethodId := ctx.Param("authMethodId")

	if err := sc.authMethodService.DeleteAuthMethod(authMethodId); err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
