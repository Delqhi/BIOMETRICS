package api

import (
	"biometrics/internal/cache"
	"biometrics/internal/config"
	"biometrics/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *cache.Redis, logger *utils.Logger, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	return gin.New()
}
