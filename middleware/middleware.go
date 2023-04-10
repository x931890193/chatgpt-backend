package middleware

import (
	"chatgpt-backend/logger"
	"chatgpt-backend/model"
	"chatgpt-backend/types"
	"chatgpt-backend/utils/useragent"
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
		// a8491c38-3b4b-4b14-8eaf-cbde20090383
		user, err := model.GetUserBySessionId(token)
		if err != nil {
			c.JSON(http.StatusOK, types.BaseResp{Message: "鉴权参数错误", Status: types.AuthError})
			c.Abort()
			return
		}
		c.Set(types.MiddlewareUser, *user)
		c.Next()
	}
}

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info.Println(useragent.ParseByRequest(c.Request))
		c.Set("client", useragent.ParseByRequest(c.Request))
		c.Next()
	}
}
