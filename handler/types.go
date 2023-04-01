package handler

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type ApiModel string

const (
	ChatGPTAPI                ApiModel = "ChatGPTAPI"
	ChatGPTUnofficialProxyAPI ApiModel = "ChatGPTUnofficialProxyAPI"
)

type MessageActionType string

const (
	Next    MessageActionType = "next"
	Variant MessageActionType = "variant"
)

type BaseResp struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
}

type ChatRequest struct {
	Options       string `json:"options"`
	Prompt        string `json:"prompt"`
	SystemMessage string `json:"systemMessage"`
}

type ChatMessage struct {
	Id              string      `json:"id"`
	Text            string      `json:"text"`
	Role            Role        `json:"role"`
	Name            string      `json:"name"`
	Delta           string      `json:"delta"`
	Detail          interface{} `json:"detail"`
	ParentMessageId string      `json:"parentMessageId"`
	ConversationId  string      `json:"conversationId"`
}

type SendMessageOptions struct {
	Name             string      `json:"name"`
	ParentMessageId  string      `json:"parentMessageId"`
	MessageId        string      `json:"messageId"`
	Stream           bool        `json:"stream"`
	SystemMessage    string      `json:"systemMessage"`
	TimeoutMs        int64       `json:"timeoutMs"`
	OnProgress       interface{} `json:"onProgress"`
	AbortSignal      interface{} `json:"abortSignal"`
	CompletionParams string      `json:"completionPparams"`
}

type SendMessageBrowserOptions struct {
	ConversationId  string            `json:"conversationId"`
	ParentMessageId string            `json:"parentMessageId"`
	MessageId       string            `json:"messageId"`
	Action          MessageActionType `json:"action"`
	TimeoutMs       int64             `json:"timeoutMs"`
	OnProgress      interface{}
}

type ChatConfig struct {
	Balance      float64 `json:"balance"`
	ReverseProxy string  `json:"reverseProxy"`
	HttpsProxy   string  `json:"httpsProxy"`
	SocksProxy   string  `json:"socksProxy"`
}
