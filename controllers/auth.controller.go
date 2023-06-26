package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/config"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
	"github.com/thuongnn/clst-mgt-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection) AuthController {
	return AuthController{authService, userService, ctx, collection}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)

	if err != nil {
		if strings.Contains(err.Error(), "username already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": newUser})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserByUsername(credentials.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not verified, please verify your email to login"})
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	appConfig, _ := config.LoadConfig(".")

	// Generate Tokens
	accessToken, err := utils.CreateToken(appConfig.AccessTokenExpiresIn, user.ID, appConfig.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(appConfig.RefreshTokenExpiresIn, user.ID, appConfig.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, appConfig.AccessTokenMaxAge*60, "/", appConfig.Domain, false, true)
	ctx.SetCookie("refresh_token", refreshToken, appConfig.RefreshTokenMaxAge*60, "/", appConfig.Domain, false, true)
	ctx.SetCookie("logged_in", "true", appConfig.AccessTokenMaxAge*60, "/", appConfig.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	appConfig, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, appConfig.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	accessToken, err := utils.CreateToken(appConfig.AccessTokenExpiresIn, user.ID, appConfig.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, appConfig.AccessTokenMaxAge*60, "/", appConfig.Domain, false, true)
	ctx.SetCookie("logged_in", "true", appConfig.AccessTokenMaxAge*60, "/", appConfig.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	appConfig, _ := config.LoadConfig(".")

	ctx.SetCookie("access_token", "", -1, "/", appConfig.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", appConfig.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", appConfig.Domain, false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
