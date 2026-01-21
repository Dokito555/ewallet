package route

import (
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	HealthController *http.HealthController
	WalletController *http.WalletController
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupPublicRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	c.App.POST("/api/healthcheck", c.HealthController.Healthcheck)
}

func (c *RouteConfig) SetupAuthRoute() {
	authGroup := c.App.Group("/api/v1/user")
	authGroup.Use(c.AuthMiddleware)
	{

	}
}
