package service

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"chatgpt-backend/types"
	"strconv"
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
	req := types.ChatRequest{
		Options:       types.ConversationRequest{},
		Prompt:        "你好呀",
		SystemMessage: "You are ChatGPT, a large language model trained by OpenAI. Follow the user's instructions carefully. Respond using markdown.",
	}
	options := &types.SendMessageBrowserOptions{
		ParentMessageId: "",
		TimeoutMs:       0,
		OnProgress:      nil,
	}
	resp, err := openAi.SendMessage(req, options)
	if err != nil {
		return
	}
	options.ParentMessageId = resp.Id
	options.ConversationId = resp.ConversationId
	req.Prompt = "我上一句话说的什么？"
	resp, err = openAi.SendMessage(req, options)
	if err != nil {
		return
	}

	options.ParentMessageId = resp.Id
	options.ConversationId = resp.ConversationId
	req.Prompt = "我之前说过什么？"
	resp, err = openAi.SendMessage(req, options)
	if err != nil {
		return
	}
	options.ParentMessageId = resp.Id
	options.ConversationId = resp.ConversationId
	req.Prompt = "从现在起你是一个科学家"
	resp, err = openAi.SendMessage(req, options)
	if err != nil {
		return
	}

	options.ParentMessageId = resp.Id
	options.ConversationId = resp.ConversationId
	req.Prompt = "作为一个科学家， 你擅长什么呢？"
	resp, err = openAi.SendMessage(req, options)
	if err != nil {
		return
	}

	options.ParentMessageId = resp.Id
	options.ConversationId = resp.ConversationId
	req.Prompt = "如何用rust实现一个http server？"
	resp, err = openAi.SendMessage(req, options)
	if err != nil {
		return
	}

	options.ConversationId = resp.ConversationId
	req.Prompt = "你现在的角色是什么？"
	resp, err = openAi.SendMessage(req, options)
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

func TestGetMessage(t *testing.T) {
	for i := 0; i < 10; i++ {
		p := strconv.Itoa(0)
		if i > 0 {
			p = strconv.Itoa(i - 1)
		} else {
			p = ""
		}
		MessageMap.Store(strconv.Itoa(i), types.ChatMessage{
			Id:              strconv.Itoa(i),
			ParentMessageId: p,
		})
	}
	message, _ := MessageMap.Load(strconv.Itoa(9))
	res := GetMessage(message.(types.ChatMessage))
	for _, message := range res {
		println(message.Id)
	}
}
