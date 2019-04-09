package routes

import "github.com/gin-gonic/gin"

const (
	ResponseError = "error"
)

func AppendResponseError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{ResponseError: err.Error()})
}
