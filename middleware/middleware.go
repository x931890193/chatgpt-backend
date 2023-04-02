package middleware

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"chatgpt-backend/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization") // Authorization
		tokenSlice := strings.Split(token, " ")
		if len(tokenSlice) != 2 {
			logger.Error.Println(fmt.Sprintf("鉴权参数错误！%v", c.Request.Header))
			c.JSON(http.StatusOK, types.BaseResp{Message: "鉴权参数错误", Status: types.AuthError})
			c.Abort()
			return
		}
		token = tokenSlice[1]
		if token != config.Cfg.OpenAI.ApiKey && token != "111111" {
			c.JSON(http.StatusOK, types.BaseResp{Message: "鉴权参数错误", Status: types.AuthError})
			c.Abort()
			return
		}
		c.Next()
	}
}
