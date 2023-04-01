package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"msg":         "chatgpt-backend",
		"server_time": time.Now().String(),
	})
}
