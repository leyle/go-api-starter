package ginhelper

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"strconv"
	"strings"
)

const DefaultCustomErrCode = 4000
const DefaultSep = "|"

type CustomErrStruct struct {
	Code int
	Msg  string
}

func (c *CustomErrStruct) Error() string {
	return fmt.Sprintf("%d|%s", c.Code, c.Msg)
}

func (c *CustomErrStruct) Append(msg string) *CustomErrStruct {
	t := &CustomErrStruct{
		Code: c.Code,
		Msg:  c.Msg + msg,
	}
	return t
}

func ParseCustomErr(err error) *CustomErrStruct {
	msg := err.Error()
	if !strings.Contains(msg, DefaultSep) {
		return &CustomErrStruct{
			Code: DefaultCustomErrCode,
			Msg:  msg,
		}
	}

	ret := strings.SplitN(msg, DefaultSep, 2)
	scode := ret[0]
	emsg := ret[1]

	return &CustomErrStruct{
		Code: parseStrCode(scode),
		Msg:  strings.TrimSpace(emsg),
	}
}

func parseStrCode(code string) int {
	code = strings.TrimSpace(code)
	c, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return DefaultCustomErrCode
	}
	return int(c)
}

func StopExec(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func RecoveryMiddleware(f func(*gin.Context, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				// runtime error, such as nil pointer dereference, should print stack
				prval := fmt.Sprintf("%v", rval)
				if strings.Contains(prval, "runtime") {
					debug.PrintStack()
				}
				err, ok := rval.(error)
				if ok {
					f(c, err)
				} else {
					err, ok := rval.(string)
					if ok {
						f(c, errors.New(err))
					} else {
						emsg := fmt.Sprintf("%v", rval)
						f(c, errors.New(emsg))
					}
				}
			}
		}()
		c.Next()
	}
}

func DefaultStopExecHandler(c *gin.Context, err error) {
	cerr := ParseCustomErr(err)
	ReturnJson(c, 400, cerr.Code, cerr.Msg, "")
}
