package handler

import "github.com/gin-gonic/gin"

func Session(c *gin.Context) {

}

func Chat(c *gin.Context) {
	req := ChatRequest{}
	if err := c.BindJSON(&req); err != nil {

		return
	}
}

func Config(c *gin.Context) {

}

func Verify(ctx *gin.Context) {

}
