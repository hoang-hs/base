package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/src/common"
	"github.com/hoang-hs/base/src/common/log"
	"net/http/httputil"
	"runtime/debug"
	"time"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				log.ErrorCtx(c.Request.Context(), "Recovery from panic", log.String("time", time.Now().String()),
					log.Any("err", err), log.String("request", string(httpRequest)), log.String("stack", string(debug.Stack())))
				e := common.ErrSystemError(c, fmt.Sprintf("recovery, err:[%s]", err))
				c.JSON(e.GetHttpStatus(), e)
				c.Abort()
			}
		}()
		c.Next()
	}
}
