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
	Model       string `gorm:"comment: 用户模型; type:VARCHAR(255)" json:"model"`
}

func (User) TableName() string {
	return "user"
}

type AccessToken struct {
	BaseModel
	SessionID  string    `gorm:"comment: 权限ID; type:VARCHAR(255); index;" json:"session_id"`
	ExpireTime time.Time `gorm:"type:timestamp; NOT NULL;" json:"-"`
	UserID     int       `gorm:"comment: 用户ID;  NOT NULL; index;" json:"user_id"`
}

func (AccessToken) TableName() string {
	return "access_token"
}

// Conversation 对话
type Conversation struct {
	BaseModel
	ConversationId string `gorm:"comment: 对话ID; index;" json:"conversation_id"`
	MessageId      string `gorm:"comment: 消息ID;  NOT NULL; index;" json:"message_id"`
	Text           string `gorm:"comment: 对话内容; NOT NULL; type:text;" json:"text"`
	AudioUrl       string `gorm:"comment: 对话语音;  NOT NULL; type:VARCHAR(255);" json:"audio_url"`
}

func (Conversation) TableName() string {
	return "conversation"
}

type Prompt struct {
	BaseModel
	Key   string `gorm:"comment: key; index; type:VARCHAR(255)" json:"key"`
	Value string `gorm:"comment: v; NOT NULL; type:text;" json:"value"`
}

// InitMysqlDB 初始化 MysqlDB 连接
func InitMysqlDB() {
	conf := config.Cfg
	password := config.Cfg.Db.Password
	host := conf.Db.Host
	if os.Getenv("PROGRAM_ENV") == "prod" {
		password = "123456"
		host = "mysql"
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
	)
	return
}