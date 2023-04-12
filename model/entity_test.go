package model

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"chatgpt-backend/service"
	"chatgpt-backend/utils"
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
)

func init() {
	InitMysqlDB()
}

func TestInitChatModel(t *testing.T) {
	models, err := service.NewAI(config.Cfg.OpenAI.ApiKey, config.Cfg.OpenAI.ApiBaseUrl).GetModels()
	if err != nil {
		logger.Error.Println(err)
		return
	}
	toSave := []GPTModel{}
	for _, model := range models.Data {
		permission := model.Permission
		marshalPermission, err := json.Marshal(permission)
		if err != nil {
			logger.Error.Println(err)
			continue
		}
		logger.Info.Println(marshalPermission)
		toSave = append(toSave, GPTModel{
			Name:       model.Id,
			Object:     model.Object,
			OwnedBy:    model.OwnedBy,
			Permission: string(marshalPermission),
		})
	}
	if err := MysqlConn.Model(&GPTModel{}).Create(toSave); err != nil {
		logger.Error.Println(err)
	}
}

func TestCreateUser(t *testing.T) {
	// a8491c38-3b4b-4b14-8eaf-cbde20090383
	u := &User{
		Avatar:      "",
		Name:        "admin" + strconv.Itoa(rand.Intn(999)),
		Description: "",
	}
	newUser, err := u.CreateUser()
	if err != nil {
		logger.Error.Println(err)
		return
	}
	token := &AccessToken{}
	accessToken, err := token.CreateToken(newUser)
	if err != nil {
		return
	}
	userModel := &UserModel{}
	model, err := userModel.CreateUserModel(newUser, 47)
	if err != nil {
		return
	}
	logger.Info.Println(accessToken.UserID, accessToken.SessionID, model.ID, model.ModelId)
}

type originPromptList struct {
	Act    string `json:"act"`
	Prompt string `json:"prompt"`
}

func TestInitPrompt(t *testing.T) {
	info, err := utils.Get("https://raw.githubusercontent.com/PlexPt/awesome-chatgpt-prompts-zh/main/prompts-zh.json", nil, utils.ContentTypeJson, nil, nil)
	if err != nil {
		logger.Error.Println(info)
		return
	}
	res := []originPromptList{}
	err = json.Unmarshal(info, &res)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	toInsert := []Prompt{}
	for _, prompt := range res {
		toInsert = append(toInsert, Prompt{
			Key:   prompt.Act,
			Value: prompt.Prompt,
		})
	}
	err = InsertPromptList(toInsert)
	if err != nil {
		logger.Error.Println(err)
		return
	}
}
