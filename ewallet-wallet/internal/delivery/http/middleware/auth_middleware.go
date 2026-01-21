package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Dokito555/ewallet/ewallet-ewallet/constants"
	"github.com/Dokito555/ewallet/ewallet-ewallet/external"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/models"
	"github.com/gin-gonic/gin"
)

// func NewAuth(userRepo *repositories.UserRepository, userService *services.UserService, tokenService *services.TokenService) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		tokenStr := ctx.GetHeader("Authorization")
// 		if tokenStr == "" {
// 			ctx.JSON(http.StatusUnauthorized, "no authorization")
// 			ctx.Abort()
// 			return
// 		}

// 		_, err := userRepo.GetUserSessionByToken(tokenStr)
// 		if err != nil {
// 			ctx.JSON(http.StatusUnauthorized, "no authorization")
// 			ctx.Abort()
// 			return
// 		}

// 		claim, err := tokenService.ValidateToken(ctx.Request.Context(), tokenStr)
// 		if err != nil {
// 			userService.Log.Warnf("failed to validate token: %+v", err)
// 			ctx.JSON(http.StatusUnauthorized, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		if time.Now().Unix() > claim.ExpiresAt.Unix() {
// 			userService.Log.Warnf("JWT token is expired. Expiry: %v, Current: %v", claim.ExpiresAt, time.Now())
// 			ctx.JSON(http.StatusUnauthorized, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("auth", claim)
// 		ctx.Next()
// 	}
// }

// func RefreshToken(userRepo *repositories.UserRepository, userService *services.UserService, tokenService *services.TokenService) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 	tokenStr := ctx.GetHeader("Authorization")
// 		if tokenStr == "" {
// 			ctx.JSON(http.StatusUnauthorized, "no authorization")
// 			ctx.Abort()
// 			return
// 		}

// 		_, err := userRepo.GetUserByRefreshToken(tokenStr)
// 		if err != nil {
// 			ctx.JSON(http.StatusUnauthorized, "no authorization")
// 			ctx.Abort()
// 			return
// 		}
// 		claim, err := tokenService.ValidateToken(ctx.Request.Context(), tokenStr)
// 		if err != nil {
// 			userService.Log.Warnf("failed to validate token: %+v", err)
// 			ctx.JSON(http.StatusUnauthorized, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		if time.Now().Unix() > claim.ExpiresAt.Unix() {
// 			userService.Log.Warnf("JWT token is expired. Expiry: %v, Current: %v", claim.ExpiresAt, time.Now())
// 			ctx.JSON(http.StatusUnauthorized, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("auth", claim)
// 		ctx.Next()
// 	}

// }


func MiddlewareValidateToken(ext *external.ExtUMS) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		tokenData, err := ext.ValidateToken(ctx.Request.Context(), auth)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		ctx.Set("auth", tokenData)

		ctx.Next()
	}
}

func MiddlewareSignatureValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientID := ctx.Request.Header.Get("Client-Id")
		if clientID == "" {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		secretKey := constants.MappingClient[clientID]
		if secretKey == "" {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		timeStamp := ctx.Request.Header.Get("Timestamp")
		if timeStamp == "" {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		requestTime, err := time.Parse(time.RFC3339, timeStamp)
		now := time.Now()
		if err != nil || now.Sub(requestTime) > 5*time.Minute {
			log.Println("invalid timestamp request")
			ctx.JSON(http.StatusInternalServerError, "no authorization",)
			ctx.Abort()
			return
		}

		signature := ctx.Request.Header.Get("Signature")
		if signature == "" {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		strPayload := ""
		if ctx.Request.Method != http.MethodGet {
			byteData, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, "no authorization")
				ctx.Abort()
				return
			}
			copyBody := io.NopCloser(bytes.NewBuffer(byteData))
			ctx.Request.Body = copyBody

			// endpoint := c.Request.URL.Path
			strPayload := string(byteData)
			re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
			strPayload = re.ReplaceAllString(strPayload, "")
			// strPayload = strings.ToLower(strPayload) + timeStamp + endpoint
			strPayload = strings.ToLower(strPayload)
		}

		h := hmac.New(sha256.New, []byte(secretKey))
		fmt.Println(strPayload)
		h.Write([]byte(strPayload))
		generatedSignature := hex.EncodeToString(h.Sum(nil))

		if signature != generatedSignature {
			ctx.JSON(http.StatusInternalServerError, "no authorization")
			ctx.Abort()
			return
		}

		ctx.Set("client_id", clientID)
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
