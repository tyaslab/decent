package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type LoggingMiddleware struct{}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (l *LoggingMiddleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		req := c.Request()
		res := c.Response()

		var requestBody bytes.Buffer
		var responseBody bytes.Buffer

		if req.Body != nil {
			teeReader := io.TeeReader(req.Body, &requestBody)
			req.Body = io.NopCloser(teeReader)
		}

		resWriter := &responseBodyWriter{
			ResponseWriter: res.Writer,
			body:           &responseBody,
		}

		res.Writer = resWriter

		err := next(c)

		duration := time.Since(start).Milliseconds()

		logEntry := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"request": map[string]interface{}{
				"method":  req.Method,
				"url":     req.URL.String(),
				"path":    c.Path(),
				"query":   req.URL.Query(),
				"headers": l.captureHeaders(req.Header),
				"body":    l.sanitizeBody(requestBody.Bytes()),
			},
			"response": map[string]interface{}{
				"status":  res.Status,
				"headers": l.captureHeaders(resWriter.Header()),
				"body":    l.sanitizeBody(responseBody.Bytes()),
			},
			"duration_ms": duration,
		}

		logJSON, _ := json.MarshalIndent(logEntry, "", "  ")
		println(string(logJSON))

		return err
	}
}

func (l *LoggingMiddleware) captureHeaders(header http.Header) map[string][]string {
	headers := make(map[string][]string)
	for k, v := range header {
		headers[k] = v
	}
	return headers
}

func (l *LoggingMiddleware) sanitizeBody(body []byte) interface{} {
	if len(body) == 0 {
		return nil
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body)
	}
	return result
}

type responseBodyWriter struct {
		http.ResponseWriter
		body *bytes.Buffer
	}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}