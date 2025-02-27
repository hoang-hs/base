package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/src/common/log"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.InfoCtx(c.Request.Context(), "received request", log.String("request", c.Request.URL.Path), log.String("method", c.Request.Method),
			log.Int("status", c.Writer.Status()), log.String("user_agent", c.Request.UserAgent()),
		)
	}
}
