package main

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	//_ "chatgpt-backend/model/conn"
	"chatgpt-backend/router"
	"context"
	"fmt"
	"net/http"
)

func main() {
	engine := router.SetupServer()
	conf := config.Cfg
	serverUrl := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	server := &http.Server{Addr: serverUrl, Handler: engine}
	logger.Info.Println(fmt.Sprintf("server listen: http://%s", serverUrl))
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err = server.Shutdown(ctx); err != nil {
		cancel()
		fmt.Println(err.Error())
	}
}
