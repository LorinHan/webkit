package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"webkit/enum"
	zhv "webkit/kit/validator"
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
	}
	ctx.JSON(http.StatusOK, resp)
}

func Fail(ctx *gin.Context, code enum.StatusCode, err error, msg ...string) {
	resp := &Resp{
		Code: code,
	}
	if err != nil {
		var vErrList validator.ValidationErrors
		if errors.As(err, &vErrList) {
			for _, vErr := range vErrList {
				if resp.Message != "" {
					resp.Message += " | "
				}
				errMsg := vErr.Translate(zhv.Trans)
				customMsg := zhv.GetCustomMsg(errMsg)
				if customMsg != "" {
					resp.Message += customMsg
				} else {
					resp.Message += errMsg
				}
			}
		} else {
			resp.Message = err.Error()
		}
	}
	if len(msg) > 0 {
		if resp.Message != "" {
			resp.Message += " | "
		}
		resp.Message += msg[0]
	}
	if resp.Message == "" {
		resp.Message = "Server exception, please try again later"
	}
	ctx.JSON(http.StatusOK, resp)
}
