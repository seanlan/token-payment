package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"token-payment/internal/api/v1"
	"token-payment/internal/middleware"
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
		paymentGroup.POST("address", v1.CreatePaymentAddress)
		paymentGroup.POST("withdraw", v1.Withdraw)
	}
	// 管理员接口
	adminGroup := apiGroup.Group("admin")
	{
		adminGroup.POST("login", v1.AdminLogin)
		adminGroup.POST("logout", middleware.RedisTokenAuthMiddleware(), v1.AdminLogout)
		adminGroup.POST("change_password", middleware.RedisTokenAuthMiddleware(), v1.AdminChangePassword)
		adminGroup.POST("info", middleware.RedisTokenAuthMiddleware(), v1.GetAdminInfo)
	}
}
