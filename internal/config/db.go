package config

import (
	"database/sql"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
	"token-payment/internal/dao"
)

type ZapWriter struct {
}

func (l ZapWriter) Printf(s string, i ...interface{}) {
	log.Infof(s, i...)
}

func NewDB(dns string) (db *gorm.DB, err error) {
	if dns == "" {
		err = errors.New("mysql dns is empty")
		return
	}
	db, err = gorm.Open(mysql.Open(dns),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "t_", SingularTable: true,
			},
			Logger: logger.New(ZapWriter{}, logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				IgnoreRecordNotFoundError: true, // 忽略记录不存在的错误
				Colorful:                  true,
			}),
		})
	if err != nil {
		return
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil || sqlDB == nil {
		return
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	return
}

func initDB(c DBConfig, debug bool) {
	db, err := NewDB(c.Uri)
	if err != nil {
		panic(err)
	}
	if debug {
		dao.DB = db.Debug()
	} else {
		dao.DB = db
	}
}
