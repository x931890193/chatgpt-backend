package microsoft

import (
	"bytes"
	"chatgpt-backend/cache"
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"chatgpt-backend/utils"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const host = "https://eastus.tts.speech.microsoft.com"
const MsToken = "MsToken"
const Expire = 10 * 50

var headers = map[string]string{
	"Ocp-Apim-Subscription-Key": config.Cfg.Microsoft.Key,
	//"Authorization": "Bearer ",
	//"Host": config.Cfg.Microsoft.Host,
}

func GetToken() string {
	token, err := cache.Client.Get("MsToken").Result()
	if err != nil || token == "" {
		body, postErr := utils.Post(config.Cfg.Microsoft.TokenUrl, nil, utils.ContentTypeJson, headers, nil)
		if postErr != nil {
			logger.Error.Printf("getToken error! %s", postErr.Error())
			return token
		}
		token = string(body)
		cache.Client.Set(MsToken, token, time.Second*Expire)
	} else {
		return token
	}
	return token
}

type voice struct {
	Name                string   `json:"Name"`
	DisplayName         string   `json:"DisplayName"`
	LocalName           string   `json:"LocalName"`
	ShortName           string   `json:"ShortName"`
	Gender              string   `json:"Gender"`
	Locale              string   `json:"Locale"`
	LocaleName          string   `json:"LocaleName"`
	StyleList           []string `json:"StyleList"`
	SampleRateHertz     string   `json:"SampleRateHertz"`
	VoiceType           string   `json:"VoiceType"`
	Status              string   `json:"Status"`
	ExtendedPropertyMap struct {
		IsHighQuality48K string `json:"IsHighQuality48K"`
	} `json:"ExtendedPropertyMap"`
	WordsPerMinute string `json:"WordsPerMinute"`
}

func GetVoiceList() []voice {
	voices := []voice{}
	body, err := utils.Get(host+"/cognitiveservices/voices/list", nil, utils.ContentTypeJson, headers, nil)
	if err != nil {
		logger.Error.Printf("get VoiceList error! %s", err.Error())
		return nil
	}
	err = json.Unmarshal(body, &voices)
	if err != nil {
		logger.Error.Printf("get VoiceList json.Unmarshal error! %s", err.Error())
		return nil
	}
	return voices
}

func TTS(text string) ([]byte, error) {
	accessToken := GetToken()
	// Set the headers
	headers := map[string]string{
		"Content-type":             "application/ssml+xml",
		"X-Microsoft-OutputFormat": "riff-24khz-16bit-mono-pcm",
		"Authorization":            "Bearer " + accessToken,
		"X-Search-AppId":           "07D3234E49CE426DAA29772419F436CA",
		"X-Search-ClientID":        "1ECFAE91408841A480F00935DC390960",
		"User-Agent":               "meidong",
	}
	xml := "<speak xmlns=\"http://www.w3.org/2001/10/synthesis\" xmlns:mstts=\"http://www.w3.org/2001/mstts\"" +
		" xmlns:emo=\"http://www.w3.org/2009/10/emotionml\" version=\"1.0\"" +
		" xml:lang=\"zh-CN\"><voice name=\"zh-CN-XiaomoNeural\"><s /><" +
		"mstts:express-as role=\"SeniorFemale\"><prosody rate=\"-20.00%\">" +
		text + "</prosody></mstts:express-as><s /></voice></speak>"

	// Send the request to the server
	logger.Info.Println("Connecting to server to synthesize the wave")
	resp, err := SendRequest("POST", "https://eastus.tts.speech.microsoft.com/cognitiveservices/v1", headers, xml)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	fmt.Println(resp.Status)
	if resp.StatusCode != 200 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}
		return nil, errors.New(string(respBody))
	}

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return respBody, nil
}

// SendRequest sends an HTTP request and returns an HTTP response
func SendRequest(method string, url string, headers map[string]string, body string) (*http.Response, error) {
	// Create an HTTP client with a 10-second timeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	// Create an HTTP request
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}

	// Set the headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Send the request and get the response
	return client.Do(req)
}
