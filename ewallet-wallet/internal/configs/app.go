package configs

import (
	"github.com/Dokito555/ewallet/ewallet-ewallet/external"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/delivery/http"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/delivery/http/middleware"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/delivery/http/route"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/repositories"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/services"
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
	walletRepo := repositories.NewWalletRepository(config.Log, config.DB)

	// ext
	umsExt := external.NewExtUMS(config.Log)

	// setup services
	walletSvc := services.NewWalletService(config.Log, config.Validate, config.Config, config.DB, walletRepo)

	// setup controllers
	walletCtrl := http.NewWalletController(config.Log, walletSvc, config.Validate)
	healthCtrl := http.NewHealthController(config.Log)

	// setup middleware
	authMiddleware := middleware.MiddlewareValidateToken(umsExt)

	// route config
	routeConfig := route.RouteConfig{
		App:              config.App,
		HealthController: healthCtrl,
		WalletController: walletCtrl,
		AuthMiddleware:   authMiddleware,
	}

	routeConfig.Setup()
}
