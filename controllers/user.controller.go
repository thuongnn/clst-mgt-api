package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.UserDBResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": models.FilteredResponse(currentUser)})
}

func (uc *UserController) FindUsers(ctx *gin.Context) {
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

	result, err := uc.userService.FindUsers(&models.UserSearchParams{
		CurrentPage:     intCurrentPage,
		PageSize:        intPageSize,
		NameKeyword:     ctx.Query("name_keyword"),
		UsernameKeyword: ctx.Query("username_keyword"),
		EmailKeyword:    ctx.Query("email_keyword"),
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

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var user *models.UserUpdate
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := uc.userService.UpdateUserById(userId, user); err != nil {
		if strings.Contains(err.Error(), "no document") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
