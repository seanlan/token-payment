package handler

import (
	"context"
	"time"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
)

func AddAdminLog(ctx context.Context, userID int64, remarks string, assetType int32, assetID int64) {
	// 添加用户日志
	_, _ = dao.AddAdminLog(ctx, &sqlmodel.AdminLog{
		AssetID:   assetID,
		AssetType: assetType,
		CreateAt:  time.Now().Unix(),
		Remarks:   remarks,
		UserID:    userID,
	})
	return
}

func AddUserMessage(ctx context.Context, userID int64, title, message, path string) {
	// 添加用户日志
	_, _ = dao.AddAdminUserMessage(ctx, &sqlmodel.AdminUserMessage{
		CreateAt: time.Now().Unix(),
		Message:  message,
		Path:     path,
		Title:    title,
		UserID:   userID,
	})
	return
}
