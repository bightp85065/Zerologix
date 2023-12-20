package dao

import (
	"context"
	"fmt"
	"zerologix/config"
	"zerologix/dao/mysql"
	"zerologix/logger"
)

var Clients = &client{}

// client ...  the dependency dao list
type client struct {
	Mysql mysql.DaoI
	// TODO kafka & redis
}

func Init(ctx context.Context) {
	Clients = &client{}
	if config.Conf.Mysql != nil {
		db, err := mysql.NewClient(ctx, config.Conf.Mysql)
		if err != nil {
			panic(fmt.Sprintf("mysql connect error %+v", err))
		}
		Clients.Mysql = db
	}
	logger.Log(ctx).Infof("Dao init successfully!")
}

func (c *client) close(ctx context.Context) {
	if c.Mysql != nil {
		c.Mysql.Close(ctx)
	}
}

func Close() {
	ctx := context.Background()
	Clients.close(ctx)
	logger.Log(ctx).Infof("Dao is closed.")
}
