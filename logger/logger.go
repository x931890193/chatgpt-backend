package logger

import (
	"chatgpt-backend/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	Error *log.Logger
	Info  *log.Logger
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func init() {
	dir := filepath.Join(config.BasePath, "log")
	exist, err := pathExists(dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if !exist {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
	fileName := filepath.Join(dir, "chatgpt-backend.log")

	setupLogger(fileName)
}

func setupLogger(fileName interface{}) {
	logFileLocation, _ := os.OpenFile(fileName.(string), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	log.SetOutput(os.Stdout)
	if os.Getenv("PROGRAM_ENV") == "prod" {
		log.SetOutput(logFileLocation)
	}
	log.SetPrefix("[ChatGpt Backend]  ")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	Error = log.New(log.Writer(), "[ChatGpt Backend] Error ", log.Flags())
	Info = log.New(log.Writer(), "[ChatGpt Backend] Info ", log.Flags())
}
