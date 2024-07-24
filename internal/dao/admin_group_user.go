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

func CountAdminGroupUser(ctx context.Context, expr clause.Expression) (totalRows int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{})
	if expr != nil {
		db = db.Where(expr)
	}
	db.Count(&totalRows)
	return totalRows, db.Error
}

func SumAdminGroupUser(ctx context.Context, sumField sqlmodel.FieldBase, expr clause.Expression) (sum float64, err error) {
	var sumValue = struct {
		N float64 `json:"n"`
	}{}
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{})
	if expr != nil {
		db = db.Where(expr)
	}
	err = db.Select("sum(" + sumField.FieldName + ") as n").Scan(&sumValue).Error
	return sumValue.N, err
}

func FetchAllAdminGroupUser(ctx context.Context, records interface{}, expr clause.Expression, page, size int, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{})
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

func FetchAdminGroupUser(ctx context.Context, record interface{}, expr clause.Expression, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{})
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

func SaveAdminGroupUser(ctx context.Context, d *sqlmodel.AdminGroupUser) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{}).Save(d)
	if err = db.Error; err != nil {
		return ErrInsertFailed
	}
	return nil
}

func AddAdminGroupUser(ctx context.Context, d *sqlmodel.AdminGroupUser) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func AddsAdminGroupUser(ctx context.Context, d *[]sqlmodel.AdminGroupUser) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func UpdateAdminGroupUser(ctx context.Context, updated *sqlmodel.AdminGroupUser) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Save(updated)
	if err = db.Error; err != nil {
		return -1, ErrUpdateFailed
	}
	return db.RowsAffected, nil
}

func UpdatesAdminGroupUser(ctx context.Context, expr clause.Expression, updated map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Updates(updated)
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}

func UpsertAdminGroupUser(ctx context.Context, d *sqlmodel.AdminGroupUser, upsert map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{}).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(upsert),
	}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func DeleteAdminGroupUser(ctx context.Context, expr clause.Expression) (rowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroupUser{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Delete(sqlmodel.AdminGroupUser{})
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}