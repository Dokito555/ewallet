package route

import (
	"github.com/Dokito555/ewallet/ewallet-ums/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	HealthController *http.HealthController
	AuthMiddleware   gin.HandlerFunc
	UserController   *http.UserController
}

func (c *RouteConfig) Setup() {
	c.SetupPublicRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	c.App.POST("/api/healthcheck", c.HealthController.Healthcheck)
	c.App.POST("/api/register", c.UserController.Register)
	c.App.POST("/api/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	authGroup := c.App.Group("/api/v1/user")
	authGroup.Use(c.AuthMiddleware)
	{
		authGroup.DELETE("/logout", c.UserController.Logout)
		authGroup.PUT("/refresh-token", c.UserController.RefreshToken)
	}
}
