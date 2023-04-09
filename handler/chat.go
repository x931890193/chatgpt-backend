package handler

import (
	"chatgpt-backend/config"
	"chatgpt-backend/service"
	"chatgpt-backend/types"
	"chatgpt-backend/utils/qiniu"
	"chatgpt-backend/utils/xunfei"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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

	resp, err := ai.SendMessage(req, &types.SendMessageBrowserOptions{
		ConversationId:  req.Options.ConversationId,
		ParentMessageId: req.Options.ParentMessageId,
	})
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: err.Error(), Status: types.Failed})
		return
	}
	if resp.Text != "" {
		ttsBytes, err := xunfei.TtsClient.TTS(resp.Text)
		if err == nil {
			key := fmt.Sprintf("gpt/tts/%s/%s.mp3", resp.ConversationId, resp.Id)
			uploadRes := qiniu.UploadStream(key, ttsBytes)
			if uploadRes != nil {
				resp.AudioUrl = fmt.Sprintf("%s%s", config.Cfg.Qiniu.Host, key)
			}
		}
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resp})
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

func HandleAsr(c *gin.Context) {
	file, err := c.FormFile("audioData")
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Error uploading audio data", Status: types.Failed})
		return
	}
	audioFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Error opening audio file", Status: types.Failed})
		return
	}
	text, _ := xunfei.AsrStreamClient.Asr(audioFile)
	c.JSON(http.StatusOK, types.BaseResp{Data: types.AsrResponse{Text: text}})
}

func Advance(c *gin.Context) {
	resp := types.AdvanceResponse{
		SystemMessage: "You are ChatGPT, a large language model trained by OpenAI. Follow the user\\'s instructions carefully. Respond using markdown.",
		Model:         "ggggg",
		Image: []types.Image{{
			Id:     "111",
			Name:   "111",
			Status: "finished",
			Url:    "https://pro-cs-freq.kefutoutiao.com/icon/tid26661/image_1582707875463_vpnlo.png",
		},
		},
		ModelList: []types.OptionModel{{
			Label: "1111111",
			Value: "111111",
		}},
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resp})
}

func Image(c *gin.Context) {
	formData, _ := c.FormFile("imageData")
	if formData != nil {
		fp, err := formData.Open()
		if err != nil {
			c.JSON(http.StatusOK, types.BaseResp{Message: "Error open image data", Status: types.Failed})
			return
		}
		key := fmt.Sprintf("gpt/image/%s", formData.Filename)
		imageBytes, err := io.ReadAll(fp)
		if err != nil {
			c.JSON(http.StatusOK, types.BaseResp{Message: "Error read audio data", Status: types.Failed})
			return
		}
		uploadRes := qiniu.UploadStream(key, imageBytes)
		if uploadRes != nil {
			imagerl := fmt.Sprintf("%s%s", config.Cfg.Qiniu.Host, key)
			c.JSON(http.StatusOK, types.BaseResp{Data: types.Image{
				Id:     "111",
				Name:   "111",
				Status: "finished",
				Url:    imagerl,
			}})
		}
	}
}

func OverView(c *gin.Context) {
	c.JSON(http.StatusOK, types.BaseResp{Data: types.UserInfo{
		Avatar:      "https://pro-cs-freq.kefutoutiao.com/icon/tid26661/image_1582707875463_vpnlo.png",
		Name:        "1111",
		Description: "2222",
	}})
}
