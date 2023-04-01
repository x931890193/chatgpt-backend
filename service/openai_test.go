package service

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"testing"
)

func TestOpenAi_GetBalance(t *testing.T) {
	openAi := OpenAi{
		ApiKey:     config.Cfg.OpenAI.ApiKey,
		ApiBaseUrl: config.Cfg.OpenAI.ApiBaseUrl,
	}
	balance, err := openAi.GetBalance()
	if err != nil {
		logger.Error.Printf("get balance error %s", err.Error())
	}
	logger.Info.Printf("balance is %f", balance)
}

func TestOpenAi_SendMessage(t *testing.T) {
	openAi := OpenAi{
		ApiKey:     config.Cfg.OpenAI.ApiKey,
		ApiBaseUrl: config.Cfg.OpenAI.ApiBaseUrl,
	}
	_, err := openAi.SendMessage("Ni好呀", nil)
	if err != nil {
		return
	}
}

func TestOpenAi_GetModels(t *testing.T) {
	openAi := OpenAi{
		ApiKey:     config.Cfg.OpenAI.ApiKey,
		ApiBaseUrl: config.Cfg.OpenAI.ApiBaseUrl,
	}
	models, err := openAi.GetModels()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	for _, model := range models.Data {
		logger.Info.Println(model.Id, model.OwnedBy)
	}
}
