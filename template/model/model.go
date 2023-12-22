package model

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
	"time"
	"webkit/config"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

func Init(dbConf config.DBConf) error {
	var (
		err     error
		conf    = &gorm.Config{}
		logConf = logger.Config{
			SlowThreshold: dbConf.SlowQueryTime,
			LogLevel:      dbConf.LogLevel,
			Colorful:      dbConf.LogColorful,
		}
	)

	if dbConf.LogLevel == 0 {
		logConf.LogLevel = logger.Warn // 默认Warn级别（慢查询）
	}
	conf.Logger = logger.New(log.New(&GormLogger{}, "", 0), logConf)

	db, err = gorm.Open(postgres.Open(dbConf.Conn), conf)
	if err != nil {
		return err
	}

	dbPool, err := db.DB()
	if err != nil {
		return err
	}

	// 设置空闲连接池中连接的最大数量
	dbPool.SetMaxIdleConns(dbConf.MaxIdleConn)

	// 设置打开数据库连接的最大数量
	dbPool.SetMaxOpenConns(dbConf.MaxOpenConn)

	// 设置连接可复用的最大时间
	if dbConf.MaxLifeTime != 0 {
		dbPool.SetConnMaxLifetime(time.Duration(dbConf.MaxLifeTime) * time.Hour)
	}
	// 设置连接处于空闲状态时可复用的最大时间
	if dbConf.MaxIdleTime != 0 {
		dbPool.SetConnMaxIdleTime(time.Duration(dbConf.MaxIdleTime) * time.Hour)
	}

	return nil
}

type GormLogger struct {
}

func (*GormLogger) Write(p []byte) (n int, err error) {
	zap.S().WithOptions(zap.WithCaller(false)).Info(strings.TrimRight(string(p), "\n"))
	return len(p), nil
}
