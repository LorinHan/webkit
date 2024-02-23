package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webkit/config"
	"webkit/kit/logger"
	"webkit/kit/validator"
	"webkit/middleware"
	"webkit/model"
	"webkit/router"
)

func main() {
	// 配置初始化
	config.InitByEnv()
	// config.InitByFile("config.yaml")

	// 日志初始化
	logger.Init(config.Conf.Logger)
	defer logger.Sync()

	engine := gin.New()
	engine.Use(middleware.Log, gin.RecoveryWithWriter(&middleware.RecoverWriter{}))
	router.Init(engine)

	// 数据库初始化
	if err := model.Init(config.Conf.DB); err != nil {
		zap.S().Fatal("数据库初始化失败", err)
	}

	// 参数校验器初始化
	if err := validator.Init(); err != nil {
		zap.S().Fatal(err)
	}

	run(engine)
}

func run(engine *gin.Engine) {
	server := &http.Server{
		Addr:    config.Conf.Server.Port,
		Handler: engine,
	}

	// 启动 HTTP 服务器（非阻塞）
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatal(err)
		}
	}()

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	// 创建一个平滑关闭的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server shutdown:", err)
	}

	zap.S().Info("Server exited")
}
