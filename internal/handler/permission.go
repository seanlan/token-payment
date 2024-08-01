package handler

import (
	"context"
	mapset "github.com/deckarep/golang-set/v2"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
)

// GetPermissionList 获取权限列表
func GetPermissionList(ctx context.Context, userId int64) (permissions []string) {
	// 查询是否是超级管理员
	var (
		userQ = sqlmodel.AdminUserColumns
		count int64
		pIDs  = make([]int, 0)
	)
	dao.DB.Model(&sqlmodel.AdminUser{}).Where(dao.And(userQ.ID.Eq(userId), userQ.IsSuper.Eq(1))).Count(&count)
	if count > 0 { // 超级管理员
		return []string{"*"}
	}
	// 查询用户所在的群组
	var (
		groupQ   = sqlmodel.AdminGroupUserColumns
		groupIDs []int
	)
	dao.DB.Model(&sqlmodel.AdminGroupUser{}).Where(groupQ.UserID.Eq(userId)).Pluck(groupQ.GroupID.FieldName, &groupIDs)
	// 查询群组所有的权限
	if len(groupIDs) > 0 {
		var (
			gpQ   = sqlmodel.AdminGroupPermissionColumns
			gpIDs []int
		)
		dao.DB.Model(&sqlmodel.AdminGroupPermission{}).Where(gpQ.GroupID.In(groupIDs)).Pluck(gpQ.PermissionID.FieldName, &gpIDs)
		if len(gpIDs) > 0 {
			pIDs = append(pIDs, gpIDs...)
		}
	}
	// 查询用户独立分配的权限
	var (
		upQ   = sqlmodel.AdminUserPermissionColumns
		upIDs []int
	)
	dao.DB.Model(&sqlmodel.AdminUserPermission{}).Where(upQ.UserID.Eq(userId)).Pluck(upQ.PermissionID.FieldName, &upIDs)
	if len(upIDs) > 0 {
		pIDs = append(pIDs, upIDs...)
	}
	// 获取权限列表
	var permissionQ = sqlmodel.AdminPermissionColumns
	dao.DB.Model(&sqlmodel.AdminPermission{}).Where(permissionQ.ID.In(pIDs)).Pluck(permissionQ.Code.FieldName, &permissions)
	return
}

// CheckPermission 权限校验
func CheckPermission(ctx context.Context, userId int64, permissions ...string) bool {
	plist := mapset.NewSet[string](GetPermissionList(ctx, userId)...)
	if plist.Contains("*") { // 全部权限
		return true
	} else {
		sublist := mapset.NewSet[string](permissions...)
		if len(plist.Intersect(sublist).ToSlice()) > 0 { // 存在交集
			return true
		}
	}
	return false
}
