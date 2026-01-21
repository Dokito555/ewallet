package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/Dokito555/ewallet/ewallet-ums/internal/models"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/repositories"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/services"
	"github.com/gin-gonic/gin"
)


func NewAuth(userRepo *repositories.UserRepository, userService *services.UserService, tokenService *services.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, "no authorization")
			ctx.Abort()
			return
		}

		_, err := userRepo.GetUserSessionByToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, "no authorization")
			ctx.Abort()
			return
		}

		claim, err := tokenService.ValidateToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			userService.Log.Warnf("failed to validate token: %+v", err)
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}

		if time.Now().Unix() > claim.ExpiresAt.Unix() {
			userService.Log.Warnf("JWT token is expired. Expiry: %v, Current: %v", claim.ExpiresAt, time.Now())
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set("auth", claim)
		ctx.Next()
	}
}

func RefreshToken(userRepo *repositories.UserRepository, userService *services.UserService, tokenService *services.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, "no authorization")
			ctx.Abort()
			return
		}

		_, err := userRepo.GetUserByRefreshToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, "no authorization")
			ctx.Abort()
			return
		}
		claim, err := tokenService.ValidateToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			userService.Log.Warnf("failed to validate token: %+v", err)
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}

		if time.Now().Unix() > claim.ExpiresAt.Unix() {
			userService.Log.Warnf("JWT token is expired. Expiry: %v, Current: %v", claim.ExpiresAt, time.Now())
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set("auth", claim)
		ctx.Next()
	}

}

func GetProfile(ctx *gin.Context) *models.ClaimToken {
	auth, exist := ctx.Get("auth")
	if !exist {
		log.Printf("auth not found in context")
		return nil
	}

	claim, ok := auth.(*models.ClaimToken)
	if !ok {
		log.Printf("type assertion failed. Expected *model.ClaimToken, got type: %T, value: %+v",
			auth, auth)
		return nil
	}

	return claim
}
