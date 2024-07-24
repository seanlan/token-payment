//generated by lazy
//author: seanlan

package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"token-payment/internal/dao/sqlmodel"
)

func CountAdminUserPermission(ctx context.Context, expr clause.Expression) (totalRows int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{})
	if expr != nil {
		db = db.Where(expr)
	}
	db.Count(&totalRows)
	return totalRows, db.Error
}

func SumAdminUserPermission(ctx context.Context, sumField sqlmodel.FieldBase, expr clause.Expression) (sum float64, err error) {
	var sumValue = struct {
		N float64 `json:"n"`
	}{}
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{})
	if expr != nil {
		db = db.Where(expr)
	}
	err = db.Select("sum(" + sumField.FieldName + ") as n").Scan(&sumValue).Error
	return sumValue.N, err
}

func FetchAllAdminUserPermission(ctx context.Context, records interface{}, expr clause.Expression, page, size int, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{})
	if expr != nil {
		db = db.Where(expr)
	}
	if page > 0 {
		db = db.Offset((page - 1) * size)
	}
	if size > 0 {
		db = db.Limit(size)
	}
	for _, order := range orders {
		db = db.Order(order)
	}
	err = db.Find(records).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
		return err
	}
	return err
}

func FetchAdminUserPermission(ctx context.Context, record interface{}, expr clause.Expression, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{})
	if expr != nil {
		db = db.Where(expr)
	}
	for _, order := range orders {
		db = db.Order(order)
	}
	err = db.First(record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
		return err
	}
	return err
}

func SaveAdminUserPermission(ctx context.Context, d *sqlmodel.AdminUserPermission) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{}).Save(d)
	if err = db.Error; err != nil {
		return ErrInsertFailed
	}
	return nil
}

func AddAdminUserPermission(ctx context.Context, d *sqlmodel.AdminUserPermission) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func AddsAdminUserPermission(ctx context.Context, d *[]sqlmodel.AdminUserPermission) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func UpdateAdminUserPermission(ctx context.Context, updated *sqlmodel.AdminUserPermission) (RowsAffected int64, err error) {
	if updated.ID == 0 {
		return -1, ErrUpdateFailed
	}
	db := GetDB(ctx).WithContext(ctx).Save(updated)
	if err = db.Error; err != nil {
		return -1, ErrUpdateFailed
	}
	return db.RowsAffected, nil
}

func UpdatesAdminUserPermission(ctx context.Context, expr clause.Expression, updated map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Updates(updated)
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}

func UpsertAdminUserPermission(ctx context.Context, d *sqlmodel.AdminUserPermission, upsert map[string]interface{}, columns ...string) (RowsAffected int64, err error) {
	var cols []clause.Column
	for _, col := range columns {
		cols = append(cols, clause.Column{Name: col})
	}
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{}).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(upsert),
		Columns:   cols,
	}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func DeleteAdminUserPermission(ctx context.Context, expr clause.Expression) (rowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminUserPermission{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Delete(sqlmodel.AdminUserPermission{})
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}
