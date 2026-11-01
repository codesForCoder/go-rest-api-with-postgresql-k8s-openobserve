package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func TraceIDResponse() gin.HandlerFunc {
	return func(c *gin.Context) {

		spanContext := trace.SpanFromContext(
			c.Request.Context(),
		).SpanContext()

		if spanContext.IsValid() {
			c.Header("trace-id", spanContext.TraceID().String())
		}

		c.Next()
	}
}
func IPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		log.Printf("IP Address: %s", ipAddress)

		// Optionally, you can store the IP address in the context
		c.Set("ClientIP", ipAddress)

		// Continue to next middleware/handler
		c.Next()
	}
}
