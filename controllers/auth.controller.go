package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thuongnn/clst-mgt-api/config"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
	"github.com/thuongnn/clst-mgt-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type AuthController struct {
	authMethodService services.AuthMethodService
	authService       services.AuthService
	userService       services.UserService
	ctx               context.Context
	collection        *mongo.Collection
}

func NewAuthController(authMethodService services.AuthMethodService, authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection) AuthController {
	return AuthController{authMethodService, authService, userService, ctx, collection}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user.Role = utils.UserRole
	user.AuthMethod = utils.BasicAuth
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

	if !user.IsActive {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Account is disabled, please contact admin for support"})
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

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	var refreshToken string

	authorizationHeader := ctx.Request.Header.Get("Authorization")
	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		refreshToken = strings.TrimPrefix(authorizationHeader, "Bearer ")
	}

	appConfig, _ := config.LoadConfig(".")
	sub, err := utils.ValidateToken(refreshToken, appConfig.RefreshTokenPublicKey)
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

	newRefreshToken, err := utils.CreateToken(appConfig.RefreshTokenExpiresIn, user.ID, appConfig.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	//currentUser := ctx.MustGet("currentUser").(*models.UserDBResponse)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) GetLoginOptions(ctx *gin.Context) {
	authMethods, err := ac.authMethodService.GetActiveAuthMethods()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Failed to get active authentication methods",
		})
		return
	}

	var response []gin.H

	for _, method := range authMethods {
		methodData := gin.H{
			"type":     method.Type,
			"name":     method.Name,
			"settings": gin.H{}, // Contains additional information if any
		}

		if method.Type == utils.Oauth2Auth {
			var oauth2Info models.OAuth2Config
			if err := json.Unmarshal(method.Configs, &oauth2Info); err != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"status":  "fail",
					"message": "Failed to parse OAuth2 config",
				})
				return
			}

			oauth2Config, err := utils.ParseOAuth2Config(oauth2Info)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
				return
			}

			authURL := oauth2Config.AuthCodeURL(method.Id.Hex(), oauth2.AccessTypeOffline)
			methodData["settings"] = gin.H{
				"auth_url":    authURL,
				"button_text": oauth2Info.ButtonText,
			}
		}

		if method.Type == utils.BasicAuth {
			var basicAuthConfig models.BasicAuthConfig
			if err := json.Unmarshal(method.Configs, &basicAuthConfig); err != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"status":  "fail",
					"message": "Failed to parse Basic Auth config",
				})
				return
			}

			methodData["settings"] = gin.H{
				"button_text": basicAuthConfig.ButtonText,
			}
		}

		response = append(response, methodData)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func (ac *AuthController) Oauth2Callback(ctx *gin.Context) {
	state, code := ctx.Query("state"), ctx.Query("code")
	if state == "" || code == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Missing state or authorization code"})
		return
	}

	authMethod, err := ac.authMethodService.GetAuthMethodById(state)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Failed to get OAuth2 config"})
		return
	}

	var oauth2Info models.OAuth2Config
	if err := json.Unmarshal(authMethod.Configs, &oauth2Info); err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Failed to parse OAuth2 config",
		})
		return
	}

	oauth2Config, err := utils.ParseOAuth2Config(oauth2Info)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var userClaims models.UserClaims
	err = utils.DecodeOauth2Token(token.AccessToken, &userClaims)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userRole := utils.UserRole
	if utils.IsAdmin(userClaims.Groups, oauth2Info.AdminGroups) {
		userRole = utils.AdminRole
	}

	userInfo, err := ac.authService.SyncOauth2User(&models.SignUpInput{
		Name:       userClaims.Name,
		Role:       userRole,
		Verified:   userClaims.EmailVerified,
		Username:   userClaims.PreferredUsername,
		Email:      userClaims.Email,
		AuthMethod: utils.Oauth2Auth,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user no logger exists"})
		return
	}

	if !userInfo.IsActive {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Account is disabled, please contact admin for support"})
		return
	}

	if !userInfo.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not verified, please verify your email to login"})
		return
	}

	// Generate Tokens
	appConfig, _ := config.LoadConfig(".")
	accessToken, err := utils.CreateToken(appConfig.AccessTokenExpiresIn, userInfo.ID, appConfig.AccessTokenPrivateKey)
	refreshToken, err := utils.CreateToken(appConfig.RefreshTokenExpiresIn, userInfo.ID, appConfig.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
