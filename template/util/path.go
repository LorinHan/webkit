package util

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
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

// GetExePath 获取当前程序执行路径
func GetExePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	return filepath.Dir(exePath) + string(os.PathSeparator)
}

// GetWorkDir 获取当前程序运行目录
func GetWorkDir() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}

	return path + string(os.PathSeparator)
}

// PathFileExists 判断文件是否存在
func PathFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsExist(err) {
		return true
	}

	return false
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// RecursiveFind 从路径向上递归查找文件，返回文件绝对路径
func RecursiveFind(path, fileName string) string {
	sp := string(os.PathSeparator)
	pathSplits := strings.Split(path+fileName, sp)

	for i := len(pathSplits); i > 0; i-- {
		filePath := strings.Join(pathSplits[0:i], sp) + sp + fileName
		if PathFileExists(filePath) {
			return filePath
		}
	}
	return ""
}

// FindConfigFile 查找配置文件
func FindConfigFile(fileName string) string {
	var path string
	path = RecursiveFind(GetExePath(), fileName)
	if path == "" {
		return RecursiveFind(GetWorkDir(), fileName)
	}
	return path
}
