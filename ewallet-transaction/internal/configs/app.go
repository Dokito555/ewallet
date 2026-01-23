package configs

import (
	"github.com/Dokito555/ewallet/ewallet-transaction/external"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/delivery/http"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/delivery/http/middleware"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/delivery/http/route"
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

	// ext
	umsExt := external.NewExtUMS(config.Log)

	// setup services
	

	// setup controllers
	healthCtrl := http.NewHealthController(config.Log)

	// setup middleware
	authMiddleware := middleware.MiddlewareValidateToken(umsExt)

	// route config
	routeConfig := route.RouteConfig{
		App:              config.App,
		HealthController: healthCtrl,
		AuthMiddleware:   authMiddleware,
	}

	routeConfig.Setup()
}
