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
	}
	{
		adminGroup.POST("logout", middleware.RedisTokenAuthMiddleware(), v1.AdminLogout)
		adminGroup.POST("change_password", middleware.RedisTokenAuthMiddleware(), v1.AdminChangePassword)
		adminGroup.POST("info", middleware.RedisTokenAuthMiddleware(), v1.GetAdminInfo)
	}
	// 使用中间件
	apiGroup.Use(middleware.RedisTokenAuthMiddleware())

	groupGroup := apiGroup.Group("group")
	{
		// 获取用户组信息
		groupGroup.POST("all", middleware.CheckPermission("group_view"), v1.GetAllGroup)
		// 新增用户组
		groupGroup.POST("add", middleware.CheckPermission("group_add"), v1.AddGroup)
		// 编辑用户组
		groupGroup.POST("edit", middleware.CheckPermission("group_edit"), v1.EditGroup)
	}
	permissionGroup := apiGroup.Group("permission")
	{
		// 获取权限列表
		permissionGroup.POST("all", v1.GetAllPermission)
	}
	staffGroup := apiGroup.Group("staff")
	{
		// 员工列表
		staffGroup.POST("list", middleware.CheckPermission("user_view"), v1.StaffList)
		// 添加员工
		staffGroup.POST("add", middleware.CheckPermission("user_add"), v1.AddStaff)
		// 编辑员工
		staffGroup.POST("edit", middleware.CheckPermission("user_edit"), v1.EditStaff)
		// 员工信息
		staffGroup.POST("info", middleware.CheckPermission("user_view"), v1.StaffInfo)
		// 重置密码
		staffGroup.POST("reset_password", middleware.CheckPermission("user_edit"), v1.StaffRePassword)
		// 员工日志
		staffGroup.POST("logs", middleware.CheckPermission("user_log_view"), v1.StaffLogs)
	}
}
