package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	http2 "social-app/pkg/http"
)

func Api() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", http2.ContentTypeJSON, "utf-8"))

		c.Next()
	}
}
