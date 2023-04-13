package handler

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"chatgpt-backend/model"
	"chatgpt-backend/service"
	"chatgpt-backend/types"
	"chatgpt-backend/utils/qiniu"
	"chatgpt-backend/utils/xunfei"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var ai = &service.OpenAi{ApiKey: config.Cfg.OpenAI.ApiKey, ApiBaseUrl: config.Cfg.OpenAI.ApiBaseUrl}

func Session(c *gin.Context) {
	c.JSON(http.StatusOK, types.BaseResp{Data: types.SessionResp{Auth: false, Model: "ChatGPTAPI"}, Status: types.Success})
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
	_, err := model.GetUserBySessionId(req.Token)
	if err != nil {
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
	for _, m := range models.Data {
		resModel = append(resModel, map[string]string{
			"label": m.Id,
			"value": m.Id,
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
	buffer, text, err := xunfei.AsrStreamClient.Asr(audioFile)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Error Parse audio data", Status: types.Failed})
		return
	}
	conversationId, _ := c.GetPostForm("conversation_Id")

	MessageId, _ := c.GetPostForm("message_Id")

	key := fmt.Sprintf("gpt/asr/%s/%s.mp3", conversationId, MessageId)

	go qiniu.UploadStream(key, buffer)

	c.JSON(http.StatusOK, types.BaseResp{Data: types.AsrResponse{Text: text}})
}

func Advance(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
	userModel, err := model.GetUserModel(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Get User Model Error!", Status: types.Failed})
		return
	}
	gptModel, err := model.GetGPTModelById(userModel.ModelId)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Get Model Info Error!", Status: types.Failed})
		return
	}
	GPTModels, err := model.GetAllGPTModels()
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Get All GPT Model Info Error!", Status: types.Failed})
		return
	}
	// "You are ChatGPT, a large language model trained by OpenAI. Follow the user's instructions carefully. Respond using markdown."
	modelList := []types.OptionModel{}
	for _, m := range GPTModels {
		modelList = append(modelList, types.OptionModel{
			Label: m.Name,
			Value: m.Name,
		})
	}
	resp := types.AdvanceResponse{
		SystemMessage: userModel.Prompt,
		Model:         gptModel.Name,
		Avatar: []types.Image{{
			Status: types.Finished,
			Url:    fmt.Sprintf("%s%s", config.Cfg.Qiniu.Host, userModel.Image),
		},
		},
		ModelList: modelList,
		Profile:   userModel.Profile,
		Name:      userModel.Name,
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resp})
}

func AdvanceSave(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
	req := types.AdvanceRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: err.Error(), Status: types.Failed})
		return
	}
	_, err := model.UpdateUserModelByUserid(user.ID, req.SystemMessage, "", req.Profile, req.Name)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Error update prompt", Status: types.Failed})
		return
	}
	c.JSON(http.StatusOK, types.BaseResp{})
}

func Image(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
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
			imageUrl := fmt.Sprintf("%s%s", config.Cfg.Qiniu.Host, key)
			_, ok = c.GetPostForm("isUser")
			if ok {
				_, err = model.UpdateUserInfoByUserid(user.ID, key, "", "")
				if err != nil {
					c.JSON(http.StatusOK, types.BaseResp{Message: "Error update image url", Status: types.Failed})
					return
				}
			} else {
				_, err = model.UpdateUserModelByUserid(user.ID, "", key, "", "")
				if err != nil {
					c.JSON(http.StatusOK, types.BaseResp{Message: "Error update image url", Status: types.Failed})
					return
				}
			}
			c.JSON(http.StatusOK, types.BaseResp{Data: types.Image{
				Status: "finished",
				Url:    imageUrl,
			}})

		} else {
			c.JSON(http.StatusOK, types.BaseResp{Message: "Error upload image data", Status: types.Failed})
		}
	}
}

func OverView(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: types.UserInfo{
		Avatar: []types.Image{{
			Status: types.Finished,
			Url:    fmt.Sprintf("%s%s", config.Cfg.Qiniu.Host, user.Avatar),
		}},
		Name:        user.Name,
		Description: user.Description,
	}})
}

func OverViewSave(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
	req := types.UserInfo{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: err.Error(), Status: types.Failed})
		return
	}
	_, err := model.UpdateUserInfoByUserid(user.ID, "", req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: "Error update General", Status: types.Failed})
		return
	}
	c.JSON(http.StatusOK, types.BaseResp{})
}

func ChatHistory(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
	history, err := model.GetChatHistoryByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: fmt.Sprintf("Get ChatHistory Error! %s", err), Status: types.Failed})
		return
	}
	resp := types.ChatHistory{}
	err = json.Unmarshal([]byte(history.History), &resp)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: fmt.Sprintf("Get ChatHistory Unmarshal Error! %s", err), Status: types.Failed})
		return
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resp})
}
func ChatHistorySave(c *gin.Context) {
	v, _ := c.Get(types.MiddlewareUser)
	user, ok := v.(model.User)
	if !ok {
		c.JSON(http.StatusOK, types.BaseResp{Message: "User Error!", Status: types.AuthError})
		return
	}
	req := types.ChatHistory{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: fmt.Sprintf("Save ChatHistorySave Error %s", err.Error()), Status: types.Failed})
		return
	}
	err := model.UpdateOrCreateChatHistory(user.ID, req)
	if err != nil {
		c.JSON(http.StatusOK, types.BaseResp{Message: fmt.Sprintf("Save ChatHistorySave Error %s", err.Error()), Status: types.Failed})
		return
	}
	c.JSON(http.StatusOK, types.BaseResp{})
}

func PromptList(c *gin.Context) {
	resp := types.PromptListResp{PromptList: []types.Prompt{}}
	promptList, err := model.GetPromptList()
	if err != nil {
		logger.Error.Println("Get prompt list Error!")
		c.JSON(http.StatusOK, types.BaseResp{Data: resp})
		return
	}
	for _, p := range promptList {
		resp.PromptList = append(resp.PromptList, types.Prompt{
			Key:   p.Key,
			Value: p.Value,
		})
	}
	c.JSON(http.StatusOK, types.BaseResp{Data: resp})

}
