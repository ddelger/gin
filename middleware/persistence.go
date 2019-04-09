package middleware

import (
	"github.com/ddelger/gin/persistence"

	"github.com/gin-gonic/gin"
)

const MiddlewarePersistence = "MiddlewarePersistence"

func Persistence(manager persistence.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(MiddlewarePersistence, manager)
		c.Next()
	}
}
