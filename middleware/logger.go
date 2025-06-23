package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ZapLogger(logger *zap.Logger, logBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate and attach request ID
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		contentType := c.ContentType()

		var body string

		if logBody && method != http.MethodGet && c.Request.Body != nil {
			switch {
			case isTextContent(contentType):
				buf, _ := io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(buf)) // Restore body
				body = string(buf)

			case strings.HasPrefix(contentType, "multipart/form-data"):
				if err := c.Request.ParseMultipartForm(10 << 20); err == nil {
					for key, headers := range c.Request.MultipartForm.File {
						for _, hdr := range headers {
							body += fmt.Sprintf("[field: %s, name: %s, size: %d, type: %s] ",
								key, hdr.Filename, hdr.Size, hdr.Header.Get("Content-Type"))
						}
					}
				}
				// NOTE: Multipart form does not need to reset c.Request.Body
			default:
				body = fmt.Sprintf("[binary content type skipped: %s]", contentType)
			}
		}

		// Call the next handler
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.String("requestID", requestID),
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.String("clientIP", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("userAgent", c.Request.UserAgent()),
		}

		if logBody && body != "" {
			fields = append(fields, zap.String("body", body))
		}

		if len(c.Errors) > 0 {
			var errList []error
			for _, e := range c.Errors {
				errList = append(errList, e.Err)
			}
			fields = append(fields, zap.Errors("errors", errList))
			logger.Error("Request encountered errors", fields...)
		} else {
			logger.Info("Handled request", fields...)
		}
	}
}

func isTextContent(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.HasPrefix(mediaType, "text/") || mediaType == "application/json" || mediaType == "application/x-www-form-urlencoded"
}
















// package middleware

// import (
// 	"bytes"
// 	"io"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"
// )

// func ZapLogger(logger *zap.Logger, logBody bool) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		path := c.Request.URL.Path
// 		raw := c.Request.URL.RawQuery

// 		// Read request body safely
// 		var body string
// 		if logBody && c.Request.Method != "GET" && c.Request.Body != nil {
// 			bodyBytes, _ := io.ReadAll(c.Request.Body)
// 			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body
// 			body = string(bodyBytes)
// 		}

// 		c.Next()

// 		latency := time.Since(start)
// 		status := c.Writer.Status()

// 		fields := []zap.Field{
// 			zap.Int("status", status),
// 			zap.String("method", c.Request.Method),
// 			zap.String("path", path),
// 			zap.String("query", raw),
// 			zap.String("clientIP", c.ClientIP()),
// 			zap.Duration("latency", latency),
// 			zap.String("userAgent", c.Request.UserAgent()),
// 		}

// 		if logBody && body != "" {
// 			fields = append(fields, zap.String("body", body))
// 		}

// 		if len(c.Errors) > 0 {
// 			var errList []error
// 			for _, e := range c.Errors {
// 				errList = append(errList, e.Err)
// 			}
// 			fields = append(fields, zap.Errors("errors", errList)) // ✅ FIXED LINE
// 			logger.Error("Request with errors", fields...)
// 		} else {
// 			logger.Info("Handled request", fields...)
// 		}
// 	}
// }
