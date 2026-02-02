package middleware

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/ginext"
)

// LoggingMiddleware логирует запросы
func LoggingMiddleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		fmt.Printf("[%s] %s %s %d %s\n",
			start.Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			status,
			time.Since(start),
		)
	}
}
