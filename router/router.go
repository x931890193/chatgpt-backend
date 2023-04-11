package router

import (
	"chatgpt-backend/handler"
	"chatgpt-backend/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

var runMode string

func init() {
	runMode = gin.DebugMode
	if os.Getenv("PROGRAM_ENV") == "prod" {
		runMode = gin.ReleaseMode
	}
}

func SetupServer() *gin.Engine {
	gin.SetMode(runMode)
	router := gin.Default()
	//router.Use(middleware.AuthMiddleware())
	//router.Use(middleware.LogMiddleware())
	router.HandleMethodNotAllowed = true
	router.GET("/", handler.Hello)
	router.Use(middleware.RequestMiddleware())
	//router.Use(middleware.BaseAuthMiddleware())
	api := router.Group("/api")
	{
		api.POST("/verify", handler.Verify)
		api.Use(middleware.AuthMiddleware())
		api.POST("/session", handler.Session)
		api.POST("/chat-process", handler.Chat)
		api.GET("/chat-history", handler.ChatHistory)
		api.POST("/config", handler.Config)
		api.POST("/chat", handler.Chat)
		api.GET("/models", handler.ModelList)
		api.POST("/audio", handler.HandleAsr)
		api.GET("/advance", handler.Advance)
		api.POST("/advance", handler.AdvanceSave)
		api.POST("/image", handler.Image)
		api.GET("/overview", handler.OverView)
		api.POST("/overview", handler.OverViewSave)

	}
	return router
}
