package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/goserve/arch/network"
)

type errorCatcher struct {
	network.BaseMiddleware
}

func NewErrorCatcher() network.RootMiddleware {
	return &errorCatcher{
		BaseMiddleware: network.NewBaseMiddleware(),
	}
}

func (m *errorCatcher) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *errorCatcher) Handler(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				m.Send(ctx).InternalServerError(err.Error(), err)
			} else {
				m.Send(ctx).InternalServerError("something went wrong", err)
			}
			ctx.Abort()
		}
	}()
	ctx.Next()
}
