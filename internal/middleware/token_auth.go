package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"token-payment/internal/constant"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/pkg/xlhttp"
)

// 根据token获取用户ID
func getUserIDByToken(ctx context.Context, token string) (userID int64, err error) {
	redisKey := fmt.Sprintf("user_token:%s", token)
	userID, err = dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		// Redis不存在
		var (
			userToken sqlmodel.AdminUserToken
			tokenQ    = sqlmodel.AdminUserTokenColumns
		)
		err = dao.FetchAdminUserToken(ctx,
			&userToken,
			dao.And(tokenQ.Token.Eq(token), tokenQ.ExpireAt.Gt(time.Now().Unix())))
		if err != nil {
			return
		}
		userID = int64(userToken.UserID)
		// 存入Redis
		_ = dao.Redis.Set(ctx, redisKey, userID, time.Second*constant.UserAccessTokenRedisExpires).Err()
	}
	return
}

func RedisTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := xlhttp.Build(c)
		var (
			req struct {
				Token string `json:"token" form:"token" binding:"omitempty"`
			}
			err    error
			userID int64
		)
		err = r.RequestParser(&req)
		if err != nil {
			c.Abort()
			return
		}
		c.Set("token", req.Token)
		userID, err = getUserIDByToken(c, req.Token)
		c.Set(xlhttp.JWTIdentityKey, userID)
		c.Next()
	}
}
