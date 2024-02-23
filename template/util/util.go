package util

import (
	"encoding/json"
	"go.uber.org/zap"
)

// Async 协程异步处理，附带recover防止panic
func Async(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error("Async Recover:", err)
			}
		}()
		f()
	}()
}

// PrintJson 将参数json格式化打印出来，方便观察、调试
func PrintJson(v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		zap.Error(err)
	}
	zap.S().WithOptions(zap.WithCaller(false)).Info(string(jsonData))
}
