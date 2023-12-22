package controller

import (
	"github.com/gin-gonic/gin"
	"webkit/enum"
	"webkit/param"
	"webkit/service"
)

var Hello = &HelloCtrl{}

type HelloCtrl struct {
}

func (*HelloCtrl) SayHi(ctx *gin.Context) {
	data, err := service.HelloSvc.SayHi()
	if err != nil {
		Fail(ctx, enum.FailedGetData, err)
		return
	}
	// Render(ctx, enum.Success, "success", data)
	Success(ctx, data)
}

func (*HelloCtrl) TestValidator(ctx *gin.Context) {
	var req param.TestValidatorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, enum.InvalidParams, err)
		return
	}

	Success(ctx, &param.TestValidatorResp{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
}
