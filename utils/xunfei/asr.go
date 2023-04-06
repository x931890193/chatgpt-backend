package xunfei

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var HOST = "wss://rtasr.xfyun.cn/v1/ws"

// END_TAG 结束标识
var END_TAG = "{\"end\": true}"

// SLICE_SIZE 每次发送的数据大小
var SLICE_SIZE = 1280

func Asr(audioContent []byte) string {
	retString := ""
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha1.New, []byte(config.Cfg.XunFei.ASRAppKey))
	strByte := []byte(config.Cfg.XunFei.AppID + ts)
	strMd5Byte := md5.Sum(strByte)
	strMd5 := fmt.Sprintf("%x", strMd5Byte)
	mac.Write([]byte(strMd5))
	signa := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	requestParam := "appid=" + config.Cfg.XunFei.AppID + "&ts=" + ts + "&signa=" + signa

	c, _, err := websocket.DefaultDialer.Dial(HOST+"?"+requestParam, nil)
	if err != nil {
		logger.Error.Println(err.Error())
		return retString
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)
	wg := &sync.WaitGroup{}
	// read
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			var result map[string]string
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			err = json.Unmarshal(message, &result)
			if err != nil {
				logger.Error.Println(string(message))
				continue
			}
			logger.Error.Println(result)
			if result["code"] == "0" {
				var asrResult AsrResult
				if result["action"] == "started" {
					logger.Info.Println(result)
					continue
				}
				err := json.Unmarshal([]byte(result["data"]), &asrResult)
				if err != nil {
					logger.Error.Println("parse asrResult error: " + err.Error())
					println("receive msg: ", string(message))
					break
				}
				if asrResult.Cn.St.Type == "0" {
					// 最终结果
					for _, wse := range asrResult.Cn.St.Rt[0].Ws {
						for _, cwe := range wse.Cw {
							retString += cwe.W
						}
					}
				} else {
					for _, wse := range asrResult.Cn.St.Rt[0].Ws {
						for _, cwe := range wse.Cw {
							print(cwe.W)
						}
					}
				}
			} else {
				println("invalid result: ", string(message))
			}

		}
	}(wg)

	// send
	wg.Add(1)
	go func(wg *sync.WaitGroup, audio []byte) {
		for i := 0; i < len(audio); i += SLICE_SIZE {
			err = c.WriteMessage(websocket.BinaryMessage, audio[:i])
			if err != nil {
				logger.Error.Println("write:", err)
				wg.Done()
			}
			//println("send data success, sleep 40 ms")
			time.Sleep(50 * time.Millisecond)
		}
		// 上传结束符
		if err := c.WriteMessage(websocket.TextMessage, []byte(END_TAG)); err != nil {
			logger.Error.Println(err.Error())
		} else {
			println("send end tag success, ", len(END_TAG))
		}
		wg.Done()
	}(wg, audioContent)

	wg.Wait()
	return retString
}

type AsrResult struct {
	Cn    Cn      `json:"cn"`
	SegId float64 `json:"seg_id"`
}

type Cn struct {
	St St `json:"st"`
}

type St struct {
	Bg   string      `json:"bg"`
	Ed   string      `json:"ed"`
	Type string      `json:"type"`
	Rt   []RtElement `json:"rt"`
}

type RtElement struct {
	Ws []WsElement `json:"ws"`
}

type WsElement struct {
	Wb float64     `json:"wb"`
	We float64     `json:"we"`
	Cw []CwElement `json:"cw"`
}

type CwElement struct {
	W  string `json:"w"`
	Wp string `json:"wp"`
}
