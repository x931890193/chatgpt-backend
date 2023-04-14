package model

import (
	"chatgpt-backend/config"
	"chatgpt-backend/logger"
	"chatgpt-backend/utils/useragent"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	MysqlConn *gorm.DB
	err       error
)

type BaseModel struct {
	ID        int       `gorm:"primary_key; index"`
	CreatedAt time.Time `gorm:"type:timestamp; NOT NULL; DEFAULT:CURRENT_TIMESTAMP;" json:"-"`
	UpdatedAt time.Time `gorm:"type:timestamp; NOT NULL; DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
	IsDelete  bool      `gorm:"default:false"`
}

// Request 各种日志
type Request struct {
	BaseModel
	IP          string `gorm:"NOT NULL; default:'127.0.0.1'" json:"ip"`
	Referer     string `gorm:"type: text" json:"referer"`
	URL         string `gorm:"NOT NULL" json:"url"`
	Major       int
	RemoteAddr  string `gorm:"NOT NULL" json:"remote_addr"`
	UserAgent   useragent.UserAgent
	OpType      string
	Method      string
	IsLogin     bool
	RequestTime uint
}

func (Request) TableName() string {
	return "request"
}

// User 用户
type User struct {
	BaseModel
	Avatar      string `gorm:"comment: 用户头像; type:VARCHAR(255)" json:"avatar"`
	Name        string `gorm:"comment: 用户名; type:VARCHAR(255)" json:"name"`
	Description string `gorm:"comment: 用户描述; type:VARCHAR(255)" json:"description"`
}

func (User) TableName() string {
	return "user"
}

// UserModel 用户对话对象
type UserModel struct {
	BaseModel
	UserID  int    `gorm:"comment: 用户ID; NOT NULL; index;" json:"user_id"`
	Image   string `gorm:"comment: 对方头像; type:VARCHAR(255)" json:"image"`
	ModelId int    `gorm:"comment: 用户模型id; index" json:"model_id"`
	Prompt  string `gorm:"comment: 模型预设信息; type:text" json:"prompt"`
	Profile string `gorm:"comment: 人物介绍; type:text" json:"profile"`
	Name    string `gorm:"comment: 人物名称; type:VARCHAR(255)" json:"name"`
}

func (UserModel) TableName() string {
	return "user_model"
}

// GPTModel 所有模型
type GPTModel struct {
	BaseModel
	Name       string `gorm:"comment: 模型信息; type:VARCHAR(255)" json:"name"`
	Object     string `gorm:"type:VARCHAR(255)" json:"object"`
	OwnedBy    string `gorm:"type:VARCHAR(255)" json:"owned_by"`
	Permission string `gorm:"type:text" json:"permission"`
}

func (GPTModel) TableName() string {
	return "gpt_model"
}

// AccessToken 登陆token
type AccessToken struct {
	BaseModel
	SessionID  string    `gorm:"comment: 权限ID; type:VARCHAR(255); index;" json:"session_id"`
	ExpireTime time.Time `gorm:"type:timestamp; NOT NULL; DEFAULT:CURRENT_TIMESTAMP;;" json:"-"`
	UserID     int       `gorm:"comment: 用户ID;  NOT NULL; index;" json:"user_id"`
}

func (AccessToken) TableName() string {
	return "access_token"
}

// Conversation 对话
type Conversation struct {
	BaseModel
	Role            string `gorm:"comment: 角色; index;" json:"role"`
	ConversationId  string `gorm:"comment: 对话ID; index;" json:"conversation_id"`
	MessageId       string `gorm:"comment: 消息ID;  NOT NULL; index;" json:"message_id"`
	ParentMessageId string `gorm:"comment: 父消息id; NOT NULL; DEFAULT: ''" json:"parent_message_id"`
	Text            string `gorm:"comment: 对话内容; NOT NULL; type:text;" json:"text"`
	AudioUrl        string `gorm:"comment: 对话语音; NOT NULL; type:VARCHAR(255);" json:"audio_url"`
}

func (Conversation) TableName() string {
	return "conversation"
}

// ConversationRelation not use
type ConversationRelation struct {
	BaseModel
	UserID         int    `gorm:"comment: 用户ID;  NOT NULL; index;" json:"user_id"`
	ConversationId string `gorm:"comment: 对话ID; index;" json:"conversation_id"`
	Active         bool   `gorm:"comment: 当前激活;  NOT NULL; index;" json:"active"`
}

// Prompt 提示词
type Prompt struct {
	BaseModel
	Key   string `gorm:"comment: key; index; type:VARCHAR(255)" json:"key"`
	Value string `gorm:"comment: v; NOT NULL; type:text;" json:"value"`
}

func (Prompt) TableName() string {
	return "prompt"
}

type ChatHistory struct {
	BaseModel
	UserID  int    `gorm:"comment: 用户ID;  NOT NULL; index;" json:"user_id"`
	History string `gorm:"comment: 所有聊天记录; type: text" json:"history"`
}

func (ChatHistory) TableName() string {
	return "chat_history"
}

// InitMysqlDB 初始化 MysqlDB 连接
func InitMysqlDB() {
	conf := config.Cfg
	password := config.Cfg.Db.Password
	host := conf.Db.Host
	if os.Getenv("PROGRAM_ENV") == "prod" {
		password = "123456"
		host = "127.0.0.1"
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Db.User,
		password,
		host,
		conf.Db.Port,
		conf.Db.Db)
	MysqlConn, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // 数据库链接配置
		SkipInitializeWithVersion: false, // 根据当前 mysql 版本自动配置
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  false, // 禁用 datetime 精度， 5.6之前不支持
		DontSupportRenameIndex:    false, // 重命名索引时采用删除并新建的方式， 5.7之前不支持
		DontSupportRenameColumn:   false, // 用 change 重命名列， 8之前和 mariadb 不支持
	}), &gorm.Config{})

	if err != nil {
		logger.Error.Println(fmt.Sprintf("connect to mysql database error: %s", err.Error()))
		panic(err.Error())
	}
	sqlDb, _ := MysqlConn.DB()
	sqlDb.SetMaxIdleConns(10)                  // 最大错误连接
	sqlDb.SetMaxOpenConns(50)                  // 最大连接数
	sqlDb.SetConnMaxLifetime(time.Second * 10) // 连接最大生命周期
	MysqlConn = MysqlConn.Debug()
	// migrate db
	_ = MysqlConn.AutoMigrate(
		&User{},
		&Request{},
		&AccessToken{},
		&Conversation{},
		&GPTModel{},
		&UserModel{},
		&Prompt{},
		&ChatHistory{},
	)
	return
}
