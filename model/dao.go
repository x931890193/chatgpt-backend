package model

import (
	"chatgpt-backend/logger"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
	"time"
)

func (user *User) CreateUser() (*User, error) {
	if err := MysqlConn.Model(user).Create(user).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return user, nil
}

func GetUserBySessionId(sessionId string) (*User, error) {
	token := &AccessToken{}
	if err := MysqlConn.Model(&AccessToken{}).Where("session_id =?", sessionId).First(token).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	user := &User{}
	if err := MysqlConn.Model(&User{}).Where("id =?", token.UserID).First(user).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return user, nil
}

func GetUserModel(userID int) (*UserModel, error) {
	u := &UserModel{}
	if err := MysqlConn.Model(&UserModel{}).Where("user_id =?", userID).First(u).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return u, nil
}

func GetGPTModelById(modelId int) (*GPTModel, error) {
	model := &GPTModel{}
	if err := MysqlConn.Model(&GPTModel{}).Where("id =?", modelId).First(model).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return model, nil
}

func GetAllGPTModels() ([]*GPTModel, error) {
	models := []*GPTModel{}
	if err := MysqlConn.Model(&GPTModel{}).Find(&models).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return models, nil
}

func UpdateUserModelByUserid(userID int, prompt, image, profile, name string) (*UserModel, error) {
	u := UserModel{}
	if prompt != "" {
		u.Prompt = prompt
	}
	if image != "" {
		u.Image = image
	}
	if profile != "" {
		u.Profile = profile
	}
	if name != "" {
		u.Name = name
	}
	if err := MysqlConn.Model(&UserModel{}).Where("user_id =?", userID).Updates(u).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return &u, nil
}

func UpdateUserInfoByUserid(userID int, avatar, name, description string) (*User, error) {
	u := User{}
	if avatar != "" {
		u.Avatar = avatar
	}
	if name != "" {
		u.Name = name
	}
	if description != "" {
		u.Description = description
	}
	if err := MysqlConn.Model(&User{}).Where("id =?", userID).Updates(u).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return &u, nil
}

func GetPromptList() ([]Prompt, error) {
	promptList := []Prompt{}
	if err := MysqlConn.Model(&Prompt{}).Find(&promptList).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return promptList, nil
}

func InsertPromptList(promptList []Prompt) error {
	db := MysqlConn.Clauses(clause.Insert{Modifier: "IGNORE"})
	if err := db.Model(&Prompt{}).CreateInBatches(&promptList, len(promptList)).Error; err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

func (t *AccessToken) CreateToken(user *User) (*AccessToken, error) {
	sessionId := uuid.New().String()
	token := &AccessToken{
		SessionID:  sessionId,
		ExpireTime: time.Now().AddDate(1, 1, 1),
		UserID:     user.ID,
	}
	if err := MysqlConn.Model(&t).FirstOrCreate(&token).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return token, nil
}

func (userModel *UserModel) CreateUserModel(user *User, ModelId int) (*UserModel, error) {
	model := &UserModel{
		UserID:  user.ID,
		Image:   "",
		ModelId: ModelId,
		Prompt:  "",
	}
	if err := MysqlConn.Model(userModel).FirstOrCreate(model).Error; err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return model, err
}
