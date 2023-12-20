package mysql

import (
	"context"
	"zerologix/config"
	"zerologix/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:generate mockery --name=DaoI --case underscore
type DaoI interface {
	BaseI
	OrderI
}

type BaseI interface {
	GetDb() *gorm.DB
	Close(ctx context.Context)
}

type dao struct {
	db *gorm.DB
}

var _ DaoI = &dao{}

func NewClient(ctx context.Context, mysqlConfig *config.MysqlConfig) (DaoI, error) {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	d := &dao{}
	d.db = db
	return d, err
}

func (d *dao) GetDb() *gorm.DB {
	return d.db
}

func (d *dao) Close(ctx context.Context) {
	sqlDB, _ := d.GetDb().DB()
	sqlDB.Close()

	// Shutdown the db
	logger.Log(ctx).Infof("DB client closed successfully")
}
