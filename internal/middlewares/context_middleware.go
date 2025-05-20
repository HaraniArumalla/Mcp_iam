package middlewares

import (
	"context"
	"iam_services_main_v1/config"

	"github.com/gin-gonic/gin"
)

// GinContextToContextMiddleware attaches the Gin context to the request context.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), config.GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
