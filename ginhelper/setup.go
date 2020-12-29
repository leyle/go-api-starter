package ginhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func SetupGin(logger *zerolog.Logger) *gin.Engine {
	e := gin.New()
	e.Use(GinLogMiddleware(logger))
	e.Use(CORSMiddleware())
	e.Use(RecoveryMiddleware(DefaultStopExecHandler))

	return e
}
