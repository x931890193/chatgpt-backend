package service

import (
	"chatgpt-backend/handler"
	"chatgpt-backend/logger"
	"chatgpt-backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type balanceResp struct {
	Data balance `json:"data"`
}

type balance struct {
	TotalAvailable float64 `json:"total_available"`
}

type ModelResp struct {
	Data []struct {
		Id         string        `json:"id"`
		Object     string        `json:"object"`
		OwnedBy    string        `json:"owned_by"`
		Permission []interface{} `json:"permission"`
	} `json:"data"`
	Object string `json:"object"`
}

type OpenAi struct {
	ApiKey     string `json:"api_key"`
	ApiBaseUrl string `json:"api_base_url"`
}

func New(apiKey, ApiBaseUrl string) *OpenAi {
	return &OpenAi{
		ApiKey:     apiKey,
		ApiBaseUrl: ApiBaseUrl,
	}
}

func (openAi *OpenAi) GetModels() (*ModelResp, error) {
	modelsUrl := fmt.Sprintf("%s/v1/models", openAi.ApiBaseUrl)
	//HTTP代理
	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)
	bytes, err := utils.Get(modelsUrl, nil, utils.ContentTypeJson, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", openAi.ApiKey),
	}, &http.Transport{
		Proxy: http.ProxyURL(proxyAddress),
	})
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}
	resp := &ModelResp{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (openAi *OpenAi) GetBalance() (float64, error) {
	balanceUrl := fmt.Sprintf("%s/dashboard/billing/credit_grants", openAi.ApiBaseUrl)
	//HTTP代理
	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)
	bytes, err := utils.Get(balanceUrl, nil, utils.ContentTypeJson, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", openAi.ApiKey),
	}, &http.Transport{
		Proxy: http.ProxyURL(proxyAddress),
	})
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}
	balance := balanceResp{}
	err = json.Unmarshal(bytes, &balance)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}
	return balance.Data.TotalAvailable, nil
}

func (openAi *OpenAi) SendMessage(text string, options *handler.SendMessageBrowserOptions) (*handler.ChatMessage, error) {
	balanceUrl := fmt.Sprintf("%s/v1/chat/completions", openAi.ApiBaseUrl)
	//HTTP代理
	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)
	reqData := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{{
			"role":    "user",
			"content": text}},
		"temperature": 0.7,
	}
	res, err := utils.Post(balanceUrl, reqData, utils.ContentTypeJson, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", openAi.ApiKey),
	}, &http.Transport{Proxy: http.ProxyURL(proxyAddress)})
	if err != nil {
		logger.Error.Println(err.Error())
	}
	logger.Info.Println(string(res))
	return nil, nil
}
