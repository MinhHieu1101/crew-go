package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cookieWriterKey struct{}

func GinCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), cookieWriterKey{}, c.Writer)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// retrieve Gin’s ResponseWriter for setting cookies
func FromContext(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(cookieWriterKey{}).(http.ResponseWriter)
	return w, ok
}
