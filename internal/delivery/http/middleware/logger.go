package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		logLevel := logrus.InfoLevel
		var levelColor, resetColor string

		switch {
		case statusCode >= 500:
			logLevel = logrus.ErrorLevel
			levelColor = "\033[31m" // Red
		case statusCode >= 400:
			logLevel = logrus.WarnLevel
			levelColor = "\033[33m" // Yellow
		default:
			levelColor = "\033[32m" // Green
		}
		resetColor = "\033[0m"

		logEntry := fmt.Sprintf("%s%s: %s [%s] %s - \"%s %s HTTP/%d.%d\" %d%s",
			levelColor,
			logLevel.String(),
			time.Now().Format("02-01-2006 15:04:05"),
			clientIP,
			userAgent,
			method,
			path,
			c.Request.ProtoMajor,
			c.Request.ProtoMinor,
			statusCode,
			resetColor,
		)

		switch logLevel {
		case logrus.ErrorLevel:
			logger.Error(logEntry)
		case logrus.WarnLevel:
			logger.Warn(logEntry)
		default:
			logger.Info(logEntry)
		}

		// Log additional details if needed
		logger.WithFields(logrus.Fields{
			"duration": duration,
			"referer":  c.Request.Referer(),
		}).Debug("Request details")
	}
}
