package middleware

import (
	"github.com/gin-gonic/gin"
	log2 "github.com/hoang-hs/base/src/common/log"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log2.InfoCtx(c.Request.Context(), "received request", log2.String("request", c.Request.URL.Path), log2.String("method", c.Request.Method),
			log2.Int("status", c.Writer.Status()), log2.String("user_agent", c.Request.UserAgent()),
		)
	}
}
