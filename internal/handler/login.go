package handler

import (
	"context"
	"fmt"
	"time"
	"token-payment/internal/constant"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/utils"
)

// GenUserToken
//
//	@Description: 生成用户token
//	@param ctx
//	@param userID
//	@return token
//	@return err
func GenUserToken(ctx context.Context, userID int64) (token sqlmodel.AdminUserToken, err error) {
	tokenQ := sqlmodel.AdminUserTokenColumns
	tokenStr := utils.GetUUIDStr()
	now := time.Now()
	expireAt := now.Unix() + constant.UserAccessTokenExpires
	refreshStr := utils.GetUUIDStr()
	refreshExpireAt := now.Unix() + constant.UserRefreshTokenExpires
	token = sqlmodel.AdminUserToken{
		UserID:          userID,
		Token:           tokenStr,
		ExpireAt:        expireAt,
		RefreshToken:    refreshStr,
		RefreshExpireAt: refreshExpireAt,
		CreateAt:        time.Now().Unix(),
	}
	_, err = dao.UpsertAdminUserToken(ctx, &token, dao.M{
		tokenQ.Token.FieldName:           tokenStr,
		tokenQ.ExpireAt.FieldName:        expireAt,
		tokenQ.RefreshToken.FieldName:    refreshStr,
		tokenQ.RefreshExpireAt.FieldName: refreshExpireAt,
		tokenQ.CreateAt.FieldName:        time.Now().Unix(),
	})
	return
}

// ClearUserToken
//
//	@Description: 清除用户token
//	@param ctx
//	@param userID
//	@return err
func ClearUserToken(ctx context.Context, userID int64) (err error) {
	var (
		tokenQ = sqlmodel.AdminUserTokenColumns
		token  sqlmodel.AdminUserToken
	)
	err = dao.FetchAdminUserToken(ctx, &token, dao.And(tokenQ.UserID.Eq(userID)))
	if err != nil {
		return
	}
	_, err = dao.DeleteAdminUserToken(ctx, dao.And(
		tokenQ.UserID.Eq(userID),
	))
	redisKey := fmt.Sprintf("user_token:%s", token.Token)
	err = dao.Redis.Del(ctx, redisKey).Err()
	return
}
