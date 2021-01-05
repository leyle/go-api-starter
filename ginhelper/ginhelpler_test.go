package ginhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

func TestGinLogMiddleware(t *testing.T) {
	e := gin.New()
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)
	e.Use(GinLogMiddleware(logger))

	e.GET("/", ginHandler)

	e.Run(":8080")
}

func ginHandler(c *gin.Context) {
	start := time.Now()
	logger := zerolog.Ctx(c.Request.Context())
	logger.Info().Msg("Start processing...")

	logger.Info().Str("id", "userapp01").Str("name", "admin").Msg("")
	//body, _ := ioutil.ReadAll(c.Request.Body)
	//logmiddleware.Debug().RawJSON("body", body).Msg("")

	c.JSON(200, gin.H{"data": start})

	logger.Info().Dur("elapsed", time.Since(start)).Msg("Done")
}

func TestExampleMain(t *testing.T) {
	ExampleMain()
}
