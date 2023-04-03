package service

import (
	"chatgpt-backend/logger"
	"chatgpt-backend/types"
	"chatgpt-backend/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"sync"
)

var MessageMap = sync.Map{}

const (
	Instructions          = "Instructions"
	UserLabelDefault      = "User"
	AssistantLabelDefault = "ChatGPT"
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

func NewAI(apiKey, ApiBaseUrl string) *OpenAi {
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

func StoreMessage(message types.ChatMessage) {
	MessageMap.Store(message.Id, message)
}

func getMessageById(messageId string) (types.ChatMessage, bool) {
	v, ok := MessageMap.Load(messageId)
	if ok {
		return v.(types.ChatMessage), true
	}
	return types.ChatMessage{}, false
}

type message struct {
	Role    types.Role `json:"role"`
	Content string     `json:"content"`
	Name    string     `json:"name"`
}

func GetMessage(msg types.ChatMessage) []types.ChatMessage {
	var messages []types.ChatMessage
	if msg.ParentMessageId != "" {
		parentMsg, ok := getMessageById(msg.ParentMessageId)
		if !ok {
			return messages
		}
		messages = append(messages, GetMessage(parentMsg)...)
	}
	messages = append(messages, msg)
	return messages
}

func (openAi *OpenAi) SendMessage(req types.ChatRequest, options *types.SendMessageBrowserOptions) (*types.ChatMessage, error) {
	balanceUrl := fmt.Sprintf("%s/v1/chat/completions", openAi.ApiBaseUrl)
	//balanceUrl := fmt.Sprintf("%s/v1/completions", openAi.ApiBaseUrl)
	//HTTP代理
	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)
	// store message
	nextMessage := types.ChatMessage{
		Id:   uuid.New().String(),
		Text: req.Prompt,
		Role: types.RoleUser,
	}

	toSaveSysTemMessage := types.ChatMessage{}
	if req.Prompt != "" {
		nextMessage.Text = req.Prompt
	}
	if options.ParentMessageId != "" {
		_, ok := getMessageById(options.ParentMessageId)
		if ok {
			nextMessage.ParentMessageId = options.ParentMessageId
		}
	} else {
		// 新建对话
		originID := uuid.New().String()
		nextMessage.ParentMessageId = originID
		toSaveSysTemMessage.Id = originID
		toSaveSysTemMessage.Role = types.RoleSystem
		if req.SystemMessage != "" {
			toSaveSysTemMessage.Text = req.SystemMessage
		}
	}
	if options.ConversationId != "" {
		nextMessage.ConversationId = options.ConversationId
	} else {
		nextMessage.ConversationId = uuid.New().String()
	}
	if toSaveSysTemMessage.Id != "" {
		StoreMessage(toSaveSysTemMessage)
	}
	historyMessage := GetMessage(nextMessage)

	toSendMessage := []message{}

	for _, msg := range historyMessage {
		toSendMessage = append(toSendMessage, message{
			Role:    msg.Role,
			Content: msg.Text,
			Name:    "ChatGPT",
		})
	}

	//systemMessageOffset := len(toSendMessage)

	//println(systemMessageOffset, toSendMessage)
	// 统计token数量
	//promptSlice := []string{}
	//for _, msg := range nextMessages {
	//	switch msg.Role {
	//	case types.RoleSystem:
	//		promptSlice = append(promptSlice, fmt.Sprintf("%s:%s", Instructions, msg.Content))
	//	case types.RoleUser:
	//		promptSlice = append(promptSlice, fmt.Sprintf("%s:%s", UserLabelDefault, msg.Content))
	//	default:
	//		promptSlice = append(promptSlice, fmt.Sprintf("%s:%s", AssistantLabelDefault, msg.Content))
	//	}
	//}
	//prompt := strings.Join(promptSlice, "\n\n")
	messageResult := types.ChatMessage{
		Detail:          types.ChatMessageDetail{},
		ParentMessageId: nextMessage.Id,
		ConversationId:  nextMessage.ConversationId,
	}
	//numTokens := 0
	// SendMessageOptions
	reqData := map[string]interface{}{
		//"model":       "gpt-3.5-turbo-0301",
		"model":       "gpt-3.5-turbo",
		"messages":    toSendMessage,
		"temperature": 0.7,
		// "top_p": number
		// "n": number
		// "max_tokens":
		// "presence_penalty": number
		// "frequency_penalty": number
		// "logit_bias": object
		// "user": string
	}
	logger.Info.Println("req data", toSendMessage)
	res, err := utils.Post(balanceUrl, reqData, utils.ContentTypeJson, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", openAi.ApiKey),
	}, &http.Transport{Proxy: http.ProxyURL(proxyAddress)})
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}

	resp := &types.BaseChatMessage{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, err
	}
	messageResult.Id = resp.Id
	if resp != nil && len(resp.Choices) > 0 {
		message2 := resp.Choices[0].Message
		messageResult.Text = message2.Content
		if message2.Role != "" {
			messageResult.Role = types.Role(message2.Role)
		}

	} else {
		messageResult.Detail = types.ChatMessageDetail{
			Choices: resp.Choices,
			Created: resp.Created,
			Id:      resp.Id,
			Model:   resp.Model,
			Object:  resp.Object,
			UseAge:  resp.Usage,
		}
	}
	StoreMessage(messageResult)
	StoreMessage(nextMessage)
	logger.Info.Println(fmt.Sprintf("role: %s | id: %s | p_id: %s | text: %s", nextMessage.Role, nextMessage.Id, nextMessage.ParentMessageId, nextMessage.Text))
	logger.Info.Println(fmt.Sprintf("role: %s | id: %s | p_id: %s | text: %s", messageResult.Role, messageResult.Id, messageResult.ParentMessageId, messageResult.Text))
	return &messageResult, nil
}
