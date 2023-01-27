package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Message []string `json:"messages"`
}

func BadRequest(c *gin.Context, messages []string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"messages": messages,
	})
}

func NotFound(c *gin.Context, messages []string) {
	c.JSON(http.StatusNotFound, gin.H{
		"messages": messages,
	})
}

func InternalServerError(c *gin.Context, messages []string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"messages": messages,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
