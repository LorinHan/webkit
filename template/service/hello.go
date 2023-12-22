package service

import (
	"go.uber.org/zap"
	"webkit/model"
)

var HelloSvc = &HelloService{}

type HelloService struct {
}

func (h *HelloService) SayHi() (string, error) {
	// return "", errors.New("test normal error")
	// return "", errork.New("test normal error")
	var count int64
	if err := model.DB().Raw("select count(id) from companies").Scan(&count).Error; err != nil {
		return "", err
	}
	zap.S().Info(count)
	return "Hello", nil
}
