package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Log(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()

	end := time.Now()
	latency := end.Sub(start)

	if len(ctx.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range ctx.Errors.Errors() {
			zap.L().Error(e)
		}
		return
	}

	zap.S().WithOptions(zap.WithCaller(false)).Info(
		"| ", ctx.Writer.Status(),
		" | ", latency,
		" | ", ctx.ClientIP(),
		" | ", ctx.Request.Method, " ",
		ctx.Request.RequestURI)
}

type RecoverWriter struct {
}

func (*RecoverWriter) Write(p []byte) (n int, err error) {
	zap.S().WithOptions(zap.WithCaller(false)).Error(string(p))
	return len(p), nil
}
