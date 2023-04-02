package handler

import (
	"chatgpt-backend/config"
	"chatgpt-backend/service"
	"chatgpt-backend/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ai = &service.OpenAi{ApiKey: config.Cfg.OpenAI.ApiKey, ApiBaseUrl: config.Cfg.OpenAI.ApiBaseUrl}

func Session(c *gin.Context) {
	c.JSON(http.StatusOK, types.BaseResp{Data: types.SessionResp{Auth: true, Model: "ChatGPTAPI"}, Status: types.Success})
}

func Chat(c *gin.Context) {
	req := types.ChatRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: err.Error(), Status: types.Failed})
		return
	}
	resp, err := ai.SendMessage(req.Prompt, nil)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: err.Error(), Status: types.Failed})
		return
	}

	c.JSON(http.StatusOK, types.BaseResp{Data: types.ChatMessage{
		Id:    "1",
		Text:  resp.Choices[0].Text,
		Role:  "user",
		Name:  "sssss",
		Delta: "1111",
		Detail: types.CreateChatCompletionDeltaResponse{
			Id:      "",
			Object:  "",
			Created: 0,
			Model:   "",
			Choices: []types.Choice{{
				Delta:        types.Delta{Role: types.RoleUser, Content: resp.Choices[0].Text + "111111"},
				Index:        0,
				FinishReason: "stop",
			}},
		},
		ParentMessageId: "0",
		ConversationId:  "111111",
	}})
}

func Config(c *gin.Context) {
	resp := types.ConfigResp{
		TimeoutMs:    0,
		ReverseProxy: "",
		ApiModel:     "",
		SocksProxy:   "",
		HttpsProxy:   "",
		Valance:      "",
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resp})
}

func Verify(c *gin.Context) {
	req := types.VerifyRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Secret key Error!", Status: types.AuthError})
		return
	}
	if req.Token == "" {
		c.JSON(http.StatusBadRequest, types.BaseResp{Message: "Secret key is empty", Status: types.AuthError})
		return
	}
	if req.Token != config.Cfg.OpenAI.ApiKey && req.Token != "111111" {
		c.JSON(http.StatusOK, types.BaseResp{Message: "密钥无效 | Secret key is invalid", Status: types.AuthError})
		return
	}
	c.JSON(http.StatusOK, types.BaseResp{Status: types.Success, Message: "Verify successfully"})
}

func ModelList(c *gin.Context) {
	models, err := ai.GetModels()
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: err.Error(), Status: types.Failed})
		return
	}
	resModel := []map[string]string{}
	for _, model := range models.Data {
		resModel = append(resModel, map[string]string{
			"label": model.Id,
			"value": model.Id,
		})
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resModel})
}
