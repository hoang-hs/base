package middleware

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func Cors() gin.HandlerFunc {
	return cors.AllowAll()
}
