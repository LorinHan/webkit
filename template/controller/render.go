package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"webkit/enum"
	zhv "webkit/kit/validator"
	"webkit/model"
)

type Resp struct {
	Code    enum.StatusCode `json:"code"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
}

func Render(ctx *gin.Context, code enum.StatusCode, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, &Resp{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func Success(ctx *gin.Context, data interface{}, msg ...string) {
	resp := &Resp{
		Code: enum.Success,
		Data: data,
	}
	if len(msg) > 0 {
		resp.Message = msg[0]
	} else {
		resp.Message = enum.Success.String()
	}
	ctx.JSON(http.StatusOK, resp)
}

func Fail(ctx *gin.Context, code enum.StatusCode, err error, msg ...string) {
	resp := &Resp{Code: code}
	if err != nil {
		var (
			vErrList validator.ValidationErrors
			dbErr    model.DBErr
		)
		if errors.As(err, &vErrList) { // 如果是validator校验的err，翻译为中文，拼接到resp.Message
			for _, vErr := range vErrList {
				if resp.Message != "" {
					resp.Message += ", "
				}
				errMsg := vErr.Translate(zhv.Trans)
				customMsg := zhv.GetCustomMsg(errMsg)
				if customMsg != "" {
					resp.Message += customMsg
				} else {
					resp.Message += errMsg
				}
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) { // 如果是gorm的record not found，翻译为中文
			resp.Message = "记录不存在"
		} else if errors.As(err, &dbErr) { // 如果是gorm其他的err，返回固定的msg，log打出实际的err信息
			zap.S().WithOptions(zap.WithCaller(false)).Error("数据库操作异常：", dbErr.Error())
			resp.Message = "数据库操作异常, 请稍后重试"
		} else {
			resp.Message = err.Error()
		}
	}

	// 如果传入msg参数，在resp.Message后面拼接上
	if len(msg) > 0 {
		if resp.Message != "" {
			resp.Message += ", "
		}
		resp.Message += msg[0]
	}

	if resp.Message == "" {
		resp.Message = "服务器异常，请稍后重试"
	}

	// 如果code描述信息不为空，在resp.Message前面拼接上
	codeMsg := code.String()
	if codeMsg != "" {
		if resp.Message != "" {
			resp.Message = ", " + resp.Message
		}
		resp.Message = codeMsg + resp.Message
	}
	ctx.JSON(http.StatusOK, resp)
}
