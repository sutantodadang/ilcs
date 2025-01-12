package middlewares

import (
	"bytes"
	"ilcs/internal/constants"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		bodyData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyData))

		log.Info().Str(constants.TRACE_ID, c.GetHeader(constants.TRACE_ID)).Str("type", "request").Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Str("query", c.Request.URL.Query().Encode()).Str("body", string(bodyData)).Send()

		c.Next()

	}

}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		responseBody := blw.body.String()
		latency := time.Since(startTime)
		log.Info().Str(constants.TRACE_ID, c.GetHeader(constants.TRACE_ID)).Str("type", "response").Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Dur("latency", latency).Str("body", responseBody).Send()
	}

}
