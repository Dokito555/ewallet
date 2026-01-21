package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HealthController struct {
	Log *logrus.Logger
}

func NewHealthController(log *logrus.Logger) *HealthController {
	return &HealthController{
		Log:  log,
	}
}

func (c *HealthController) Healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "healthy",
	})
}