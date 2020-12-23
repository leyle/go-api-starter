package ginhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// below is an example
type ExampleContext struct {
	C   *gin.Context
	Log *zerolog.Logger
	// more fields
}

func (a *ExampleContext) New(c *gin.Context) *ExampleContext {
	logger := log.Ctx(c.Request.Context())
	ac := &ExampleContext{
		C:   c,
		Log: logger,
	}
	return ac
}

func HandlerWrapper(f func(ctx *ExampleContext), ctx *ExampleContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		nctx := ctx.New(c)
		f(nctx)
	}
}

func exampleHandler(ctx *ExampleContext) {
	ctx.Log.Info().Str("user", "user0001").Msg("start process request")
	ctx.C.JSON(200, "OK")
	ctx.Log.Info().Str("type", "end").Msg("")
}

func ExampleMain() {
	e := gin.New()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	e.Use(GinLogMiddleware(logger))

	ctx := &ExampleContext{}
	e.GET("/", HandlerWrapper(exampleHandler, ctx))
	e.Run(":8080")
}
