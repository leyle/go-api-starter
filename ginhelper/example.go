package ginhelper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/httpclient"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
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
	time.Sleep(5 * time.Second)
	ctx.C.JSON(200, "OK")
	ctx.Log.Info().Str("type", "last").Msg("")
}

func httpClientHandler(ctx *ExampleContext) {
	// test gin middleware and httpclient api
	url := "http://192.168.2.40:8000/get"
	query := make(map[string][]string)
	query["name"] = []string{"jack", "telsa"}
	query["age"] = []string{"12", "23"}

	headers := make(map[string]string)
	headers["X-REQ-ID"] = "hello world"

	cReq := &httpclient.ClientRequest{
		Ctx:     ctx.C.Request.Context(),
		Url:     url,
		Query:   query,
		Headers: headers,
		Timeout: 10,
		Debug:   true,
	}

	resp := httpclient.Get(cReq)
	if resp.Err != nil {
		ctx.C.JSON(400, resp.Err.Error())
		return
	}

	ctx.C.JSON(200, string(resp.Body))
	return
}

func panicHandler(ctx *ExampleContext) {
	err := errors.New("should stop")
	StopExec(err)
}

func ExampleMain() {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)

	e := SetupGin(logger)

	ctx := &ExampleContext{}
	e.GET("/", HandlerWrapper(exampleHandler, ctx))
	e.GET("/get", HandlerWrapper(httpClientHandler, ctx))
	e.GET("/panic", HandlerWrapper(panicHandler, ctx))
	e.Run(":8080")
}
