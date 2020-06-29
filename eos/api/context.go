package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data interface{}) {
	result := gin.H{"success": true, "data": data}

	c.JSON(http.StatusOK, result)
}
