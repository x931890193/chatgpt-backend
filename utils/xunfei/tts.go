package xunfei

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	"time"
)

var urlBase = "wss://tts-api.xfyun.cn/v2/tts"

var TtsClient = &TTS{
	AppKey:    config.Cfg.XunFei.TTSAppKey,
	AppID:     config.Cfg.XunFei.AppID,
	AppSecret: config.Cfg.XunFei.TTSAppSecret,
}

type TTS struct {
	AppKey    string
	AppID     string
	AppSecret string
}

type TTSResultData struct {
	Audio  string `json:"audio"`
	Status int    `json:"status"`
	Ced    string `json:"ced"`
}

type TTSResult struct {
	Data    TTSResultData `json:"data"`
	Code    int
	Message string
}

type businessArgs struct {
	Aue   string `json:"aue"`
	Auf   string `json:"auf"`
	Vcn   string `json:"vcn"`
	Tte   string `json:"tte"`
	Sfl   int    `json:"sfl"`
	Speed int    `json:"speed"`
}

type common struct {
	AppId string `json:"app_id"`
}
type data struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

type ttsRequest struct {
	Common   common       `json:"common"`
	Business businessArgs `json:"business"`
	Data     data         `json:"data"`
}

type ttsSendEnd struct {
	Data interface{} `json:"data"`
}

func (tts *TTS) createUrl() string {

	// Generate RFC1123 formatted timestamp
	date := time.Now().UTC().Format(time.RFC1123)

	// Concatenate signature origin string
	signatureOrigin := "host: ws-api.xfyun.cn\n"
	signatureOrigin += "date: " + date + "\n"
	signatureOrigin += "GET /v2/tts HTTP/1.1"

	// Compute signature
	mac := hmac.New(sha256.New, []byte(tts.AppSecret))
	mac.Write([]byte(signatureOrigin))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Concatenate authorization origin string
	authorizationOrigin := fmt.Sprintf("api_key=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"",
		tts.AppKey, "hmac-sha256", "host date request-line", signature)

	// Encode authorization string with base64
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))

	// Create URL query parameters
	v := url.Values{}
	v.Set("authorization", authorization)
	v.Set("date", date)
	v.Set("host", "ws-api.xfyun.cn")

	// Concatenate query parameters to URL
	urlBase += "?" + v.Encode()

	return urlBase
}

func (tts *TTS) newTTSRequest(text string) *ttsRequest {
	return &ttsRequest{
		Common: common{AppId: tts.AppID},
		Business: businessArgs{
			Aue:   "lame",
			Auf:   "audio/L16;rate=16000",
			Vcn:   "xiaoyan",
			Tte:   "utf8",
			Sfl:   1,
			Speed: 10,
		},
		Data: data{
			Status: 2,
			Text:   base64.StdEncoding.EncodeToString([]byte(text)),
		},
	}
}

func (tts *TTS) TTS(text string) []byte {
	resStream := []byte{}
	wssUrl := tts.createUrl()
	c, _, err := websocket.DefaultDialer.Dial(wssUrl, nil)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)
	req := tts.newTTSRequest(text)
	reqByte, _ := json.Marshal(req)
	wg := &sync.WaitGroup{}
	// send
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			err := c.WriteMessage(websocket.TextMessage, reqByte)
			if err != nil {
				logger.Error.Println(fmt.Sprintf("send msg error! %s", reqByte))
				break
			}
			break
		}
	}(wg)
	// read
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			result := TTSResult{}
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Error.Println("read:", err)
				break
			}
			err = json.Unmarshal(message, &result)
			if err != nil {
				logger.Error.Println(err.Error())
				break
			}
			if result.Code != 0 {
				logger.Error.Println(result.Message)
			} else {
				res, err := base64.StdEncoding.DecodeString(result.Data.Audio)
				if err != nil {
					logger.Error.Println("base64 decode Error")
					return
				}
				resStream = append(resStream, res...)
			}
			if result.Data.Status == 2 {
				logger.Error.Println("ws is closed")
				break
			}
		}
	}(wg)
	wg.Wait()
	return resStream
}
