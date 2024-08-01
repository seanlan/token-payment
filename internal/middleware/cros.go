package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const CrossRequestHeader = "XXX-CROSS-DISABLE"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		crossHeader := c.Request.Header.Get(CrossRequestHeader)
		if crossHeader == "" { // 如果没有设置跨域请求头，则默认允许所有跨域请求
			//origin := c.Request.Header.Get("Origin") //请求头部
			accessControlHeader := c.Request.Header.Get("Access-Control-Request-Headers")
			//if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "*")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "DNTX-Mx-Red Token,Keep-Alive,User-Agent.X-Requested-With lf-Modified-Since,Cache-Control,Content-Type,Authorization")
			c.Header("Access-Control-Allow-Headers", accessControlHeader)
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
			//}
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
