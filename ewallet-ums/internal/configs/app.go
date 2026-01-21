package configs

import (
	"github.com/Dokito555/ewallet/ewallet-ums/internal/delivery/http"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/delivery/http/middleware"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/delivery/http/route"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/repositories"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repo
	userRepo := repositories.NewUserRepository(config.Log, config.DB)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config, userRepo)
	userService := services.NewUserService(config.Log, config.Validate, config.Config, config.DB, userRepo, tokenService)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	userController := http.NewUserController(config.Log, userService, config.Validate)

	// setup middleware
	authMiddleware := middleware.NewAuth(userRepo, userService, tokenService)
	// refreshToken := middleware.RefreshToken(userRepo, userService, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App:              config.App,
		HealthController: healthController,
		AuthMiddleware:   authMiddleware,
		UserController:   userController,
	}

	routeConfig.Setup()
}
