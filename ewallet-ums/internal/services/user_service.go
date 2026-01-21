package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokito555/ewallet/ewallet-ums/internal/entity"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/models"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	DB       *gorm.DB
	UserRepo *repositories.UserRepository
	TokenService *TokenService
}

func NewUserService(log *logrus.Logger, val *validator.Validate, config *viper.Viper, db *gorm.DB, userRepo *repositories.UserRepository, tokenService *TokenService) *UserService {
	return &UserService{
		Log:      log,
		Validate: val,
		Config:   config,
		DB:       db,
		UserRepo: userRepo,
		TokenService: tokenService,
	}
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {

	userDetail, err := s.UserRepo.GetUserByUsername(req.Username)
	if err != nil {
		s.Log.Warnf("failed to get user by username: %v", err)
		return nil, fmt.Errorf("failed to get user by username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(req.Password)); err != nil {
		s.Log.Warnf("failed to compare password: %v", err)
		return nil, fmt.Errorf("failed to compare password")
	}

	token, err := s.TokenService.GenerateToken(ctx, userDetail.ID, userDetail.Username, userDetail.FullName, "token", userDetail.Email)
	if err != nil {
		s.Log.Warnf("failed to generate password")
		return nil, fmt.Errorf("failed to generate token")
	}

	refreshToken, err := s.TokenService.GenerateToken(ctx, userDetail.ID, userDetail.Username, userDetail.FullName, "token", userDetail.Email)
	if err != nil {
		s.Log.Warnf("failed to generate refresh token")
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	userSession := &entity.UserSession{
		UserID:              int(userDetail.ID),
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        time.Now().Add(time.Hour * 24),
		RefreshTokenExpired: time.Now().Add(time.Hour * 72),
	}

	err = s.UserRepo.InsertNewUserSession(*userSession)
	if err != nil {
		s.Log.Warnf("failed to insert new user")
		return nil, fmt.Errorf("failed to insert new user")
	}

	rsp := &models.LoginResponse{
		UserID: userDetail.ID,
		Username: userDetail.Username,
		FullName: userDetail.FullName,
		Email: userDetail.Email,
		Token: token,
		RefreshToken: refreshToken,
	}

	return rsp, nil
}

func (s *UserService) Logout(ctx context.Context, token string) (error) {
	return s.UserRepo.DeleteUserSession(token)
} 

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string, claim models.ClaimToken) (*models.RefreshTokenResponse, error) {
	token, err := s.TokenService.GenerateToken(ctx, claim.UserID, claim.Username, claim.FullName, "token", claim.Email)
	if err != nil {
		s.Log.Warnf("failed to generate token")
		return nil, fmt.Errorf("failed to generate token")
	}

	err = s.UserRepo.UpdateTokenByRefreshToken(token, refreshToken)
	if err != nil {
		s.Log.Warnf("failed to update new token")
		return nil, fmt.Errorf("failed to update new token")
	}

	rsp := &models.RefreshTokenResponse{
		Token: token,
	}

	return rsp, nil
} 

func (s *UserService) Register(ctx context.Context, request entity.User) (interface{}, error) {

	user, err := s.UserRepo.GetUserByEmail(request.Email)
	if user.Email == request.Email {
		s.Log.Warnf("user with that email already existed")
		return nil, fmt.Errorf("user with that email already existed")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	request.Password = string(hashPassword)

	err = s.UserRepo.Create(s.DB, &request)
	if err != nil {
		return nil, err
	}

	// _, err = s.External.CreateWallet(ctx, request.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// err = s.External.SendNotification(ctx, request.Email, "register", map[string]string{
	// 	"full_name": request.FullName,
	// })

	rsp := request
	rsp.Password = ""

	return rsp, nil
}