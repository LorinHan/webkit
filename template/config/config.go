package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"os"
	"time"
	"webkit/util"
)

/**
 * config，配置包
   - 加载配置，两种方式：
     InitByEnv-通过环境变量加载配置
	 InitByFile-通过读取文件加载配置（支持yaml、json、toml等多种格式）
   - 取用配置：config.Conf.XX
*/

var Conf Config

// InitByEnv 通过环境变量加载配置
func InitByEnv() {
	Conf.DB = DBConf{
		Conn: GetEnvString("DB_CONN", ""),
	}
}

// InitByFile 通过文件加载配置
func InitByFile(fileName string) {
	viper.SetConfigFile(util.FindConfigFile(fileName))
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Panic("config init fail", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		zap.S().Panic("config init fail", err)
	}
	viper.Reset() // 释放viper内存
}

// GetEnvString 获取环境变量
func GetEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) != 0 {
		return value
	}
	return defaultValue
}

type Config struct {
	DB DBConf `json:"db"`
}

type DBConf struct {
	Conn          string          `json:"conn"`
	MaxIdleConn   int             `json:"max_idle_conn"`   // 最大空闲连接
	MaxOpenConn   int             `json:"max_open_conn"`   // 最大连接数
	MaxLifeTime   int             `json:"max_life_time"`   // 最大活跃时间，单位：h
	MaxIdleTime   int             `json:"max_idle_time"`   // 最大空闲保活时间，单位：h
	SlowQueryTime time.Duration   `json:"slow_query_time"` // 慢 SQL 阈值
	LogLevel      logger.LogLevel `json:"log_level"`       // 日志等级
	LogColorful   bool            `json:"log_colorful"`    // 启用彩色日志
}
