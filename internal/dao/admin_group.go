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

func CountAdminGroup(ctx context.Context, expr clause.Expression) (totalRows int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{})
	if expr != nil {
		db = db.Where(expr)
	}
	db.Count(&totalRows)
	return totalRows, db.Error
}

func SumAdminGroup(ctx context.Context, sumField sqlmodel.FieldBase, expr clause.Expression) (sum float64, err error) {
	var sumValue = struct {
		N float64 `json:"n"`
	}{}
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{})
	if expr != nil {
		db = db.Where(expr)
	}
	err = db.Select("sum(" + sumField.FieldName + ") as n").Scan(&sumValue).Error
	return sumValue.N, err
}

func FetchAllAdminGroup(ctx context.Context, records interface{}, expr clause.Expression, page, size int, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{})
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

func FetchAdminGroup(ctx context.Context, record interface{}, expr clause.Expression, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{})
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

func SaveAdminGroup(ctx context.Context, d *sqlmodel.AdminGroup) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{}).Save(d)
	if err = db.Error; err != nil {
		return ErrInsertFailed
	}
	return nil
}

func AddAdminGroup(ctx context.Context, d *sqlmodel.AdminGroup) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func AddsAdminGroup(ctx context.Context, d *[]sqlmodel.AdminGroup) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func UpdateAdminGroup(ctx context.Context, updated *sqlmodel.AdminGroup) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Save(updated)
	if err = db.Error; err != nil {
		return -1, ErrUpdateFailed
	}
	return db.RowsAffected, nil
}

func UpdatesAdminGroup(ctx context.Context, expr clause.Expression, updated map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Updates(updated)
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}

func UpsertAdminGroup(ctx context.Context, d *sqlmodel.AdminGroup, upsert map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{}).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(upsert),
	}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func DeleteAdminGroup(ctx context.Context, expr clause.Expression) (rowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.AdminGroup{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Delete(sqlmodel.AdminGroup{})
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}
