package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Err(e)
			}
		} else {
			l := logger.Info().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("ip", c.ClientIP()).
				Str("user-agent", c.Request.UserAgent()).
				Int64("timestamp", end.Unix()).
				Time("time", end).
				Str("latency", fmt.Sprintf("%v", latency)).
				Str("path", path)

			if query != "" {
				l.Str("query", query)
			}
			l.Msg("REQUEST")
		}
	}
}
