package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Log struct{}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		c.Next()

		status := c.Writer.Status()
		log.Printf("Request: %s %s from %s, Response Status: %d", method, path, clientIP, status)
	}
}
