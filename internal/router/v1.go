package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"token-payment/pkg/gin_zap"
)

func initWebRouter(r *gin.Engine) {
	r.Use(
		gin_zap.Ginzap(zap.L(), time.RFC3339, false),
		gin_zap.RecoveryWithZap(zap.L(), true),
	)

	// 接口处理
	apiGroup := r.Group("api/v1")
	paymentGroup := apiGroup.Group("payment")
	{

	}
}
