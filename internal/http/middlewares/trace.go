package middlewares

import (
	"ilcs/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {

		trace := c.GetHeader(constants.TRACE_ID)

		if trace == "" {
			c.Request.Header.Add(constants.TRACE_ID, xid.New().String())
		}

		c.Next()

	}
}
