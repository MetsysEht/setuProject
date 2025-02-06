package server

import (
	"github.com/MetsysEht/setuProject/pkg/logger"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

var S *gin.Engine

func Initialize() {
	gin.SetMode(gin.ReleaseMode)
	S = gin.New()
	S.Use(cors.Default())
	//S.Use(middleware.CheckAuthMiddleware)
	S.Use(ginzap.RecoveryWithZap(logger.L.Desugar(), true))

	registerRoutes()
}

func registerRoutes() {

}
