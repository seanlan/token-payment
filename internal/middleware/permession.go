package middleware

import (
	"github.com/gin-gonic/gin"
	"token-payment/internal/e"
	"token-payment/internal/handler"
	"token-payment/pkg/xlhttp"
)

// CheckPermission 校验权限
func CheckPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := xlhttp.Build(c)
		userID, err := r.GetJWTUID()
		if userID == 0 || err != nil {
			r.JsonReturn(e.ErrorToken)
			c.Abort()
			return
		}
		if handler.CheckPermission(c, userID, permissions...) {
			c.Next()
		} else {
			r.JsonReturn(e.ErrorPermission)
			c.Abort()
		}
		c.Next()
	}
}
