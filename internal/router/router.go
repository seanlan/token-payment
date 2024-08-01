package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"token-payment/internal/middleware"
)

// Setup 初始化Router
func initRouter(debug bool) *gin.Engine {
	//设置启动模式
	switch debug {
	case true:
		gin.SetMode(gin.DebugMode)
	case false:
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Cors())
	return r
}

func Run(addr string, debug bool) {
	r := initRouter(debug)
	initWebRouter(r)
	err := r.Run(addr)
	if err != nil {
		zap.S().Fatalf("run web server error: %v", err)
	}
}
