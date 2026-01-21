package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokito555/ewallet/ewallet-ums/internal/models"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/repositories"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TokenService struct {
	Log   *logrus.Logger
	Viper *viper.Viper
	UserRepo *repositories.UserRepository
}

func NewTokenService(log *logrus.Logger, viper *viper.Viper, userRepo *repositories.UserRepository) *TokenService {
	return &TokenService{
		Log:   log,
		Viper: viper,
		UserRepo: userRepo,
	}
}

var mapTypeToken = map[string]time.Duration{
	"token": time.Hour*24,
	"refresh_token": time.Hour*72,
}

func (s *TokenService) GenerateToken(ctx context.Context, userID int, username, fullname, tokenType string, email string) (string, error) {
	claimToken := models.ClaimToken{
		UserID: userID,
		Username: username,
		FullName: fullname,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: 	s.Viper.GetString("APP_NAME"),
			IssuedAt: 	jwt.NewNumericDate(time.Now()),
			ExpiresAt: 	jwt.NewNumericDate(time.Now().Add(mapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	resultToken, err := token.SignedString(s.Viper.GetString("APP_SECRET"))
	if err != nil {
		return resultToken, fmt.Errorf("failed to generate token: %v", err)
	}
	return resultToken, err
}

func (s *TokenService) ValidateToken(ctx context.Context, token string) (*models.ClaimToken, error) {
	var (
		claimToken *models.ClaimToken
		ok bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &models.ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to validate method jwt: %v", t.Header["alg"])
		}
		return s.Viper.GetString("APP_SECRET"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*models.ClaimToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}

func (s *TokenService) TokenValidation(ctx context.Context, token string) (*models.ClaimToken, error) {
	var (
		claimToken *models.ClaimToken
		err error
	)

	claimToken, err = s.ValidateToken(ctx, token)
	if err != nil {
		s.Log.Warnf("failed to validate token")
		return claimToken, fmt.Errorf("failed to validate token")
	}

	_, err = s.UserRepo.GetUserSessionByToken(token)
	if err != nil {
		s.Log.Warnf("failed to get user session")
		return claimToken, fmt.Errorf("failed to get user session")
	}
	
	return claimToken, nil
}