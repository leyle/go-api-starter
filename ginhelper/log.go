package ginhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/rs/zerolog"
	"io/ioutil"
	"time"
)

var PrintHeaders = false

var ignoreReadReqBodyPath = []string{}

func AddIgnoreReadReqBodyPath(paths ...string) {
	ignoreReadReqBodyPath = append(ignoreReadReqBodyPath, paths...)
}

func isIgnoreReadBodyPath(reqPath string) bool {
	for _, path := range ignoreReadReqBodyPath {
		if reqPath == path {
			return true
		}
	}
	return false
}

// rewrite Write()
type respWriter struct {
	gin.ResponseWriter
	cache *bytes.Buffer
}

// it will increase memory usage
func (r *respWriter) Write(b []byte) (int, error) {
	r.cache.Write(b)
	return r.ResponseWriter.Write(b)
}

func GinLogMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startT := time.Now()
		id := c.Request.Header.Get(logmiddleware.ReqIdHeaderName)
		if id == "" {
			id = logmiddleware.GenerateReqId()
		}

		// save logger into current gin.Context
		// setup logger req id and save it into current gin.Context
		// can be used in later code
		// e.g. logger := zerolog.Ctx(c.Request.Context())
		l := logger.With().Str(logmiddleware.ReqIdContextName, id).Logger()
		lctx := l.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(lctx)

		// save req id into current gin.Context
		// can be used in later code
		// e.g. c.Get(logmiddleware.ReqIdContextName)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, logmiddleware.ReqIdContextName, id)
		c.Request = c.Request.WithContext(ctx)

		// print req logmiddleware
		body := reqBody(c)
		reqInfo := reqJson(c)
		event := l.Debug().Str("type", "REQUEST").RawJSON("req", reqInfo).RawJSON("body", body)
		if PrintHeaders {
			headers, _ := json.Marshal(c.Request.Header)
			event.RawJSON("headers", headers)
		}
		event.Msg("")

		// write req id to response headers
		c.Writer.Header().Set(logmiddleware.ReqIdHeaderName, id)

		c.Writer = &respWriter{
			ResponseWriter: c.Writer,
			cache:          bytes.NewBufferString(""),
		}

		c.Next()

		// write response info into logmiddleware
		statusCode := c.Writer.Status()
		revent := l.Info().Str("type", "RESPONSE").Int("statusCode", statusCode)
		rw, ok := c.Writer.(*respWriter)
		if !ok {
			// silently passed
		} else {
			if rw.cache.Len() > 0 && !isIgnoreReadBodyPath(c.Request.URL.Path) {
				revent.RawJSON("body", rw.cache.Bytes())
			}
		}
		latency := time.Since(startT)
		revent.Str("latency", latency.String()).Msg("")
	}
}

func reqJson(c *gin.Context) []byte {
	path := c.Request.RequestURI
	method := c.Request.Method
	ctype := c.Request.Header.Get("Content-Type")
	clientIp := c.ClientIP()

	data := map[string]string{
		"path":   path,
		"method": method,
		"ctype":  ctype,
		"ip":     clientIp,
	}
	bdata, _ := json.Marshal(data)
	return bdata
}

func reqBody(c *gin.Context) []byte {
	var err error
	var body []byte
	if c.Request.ContentLength > 0 && !isIgnoreReadBodyPath(c.Request.URL.Path) {
		body, err = ioutil.ReadAll(c.Request.Body)
		if err == nil {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
	}
	return body
}
