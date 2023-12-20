package mysql

import (
	"context"

	mysqlModel "zerologix/model/mysql"

	"github.com/pkg/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderI interface {
	OrderCreate(ctx context.Context, tx *gorm.DB, row *mysqlModel.Order) error
}

func (d *dao) OrderCreate(ctx context.Context, tx *gorm.DB, row *mysqlModel.Order) (err error) {
	if tx == nil {
		tx = d.GetDb()
	}

	tx = tx.Table(mysqlModel.Order{}.TableName())
	if err = tx.Clauses(clause.Returning{}).Create(row).Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}
