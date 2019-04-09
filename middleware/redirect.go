package middleware

import (
	"net/http"

	"github.com/ddelger/glog"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func Redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := secure.New(secure.Options{SSLRedirect: true}).Process(c.Writer, c.Request); err != nil {
			glog.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}

		c.Next()
	}
}
