package types

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type ApiModel string

const (
	MiddlewareUser                     = "MiddlewareUser"
	Finished                           = "finished"
	ChatGPTAPI                ApiModel = "ChatGPTAPI"
	ChatGPTUnofficialProxyAPI ApiModel = "ChatGPTUnofficialProxyAPI"
)

type MessageActionType string

const (
	Next    MessageActionType = "next"
	Variant MessageActionType = "variant"
)

type RespCode int

const (
	Success   RespCode = 0
	Failed    RespCode = 1
	AuthError RespCode = 2
)

type BaseResp struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  RespCode    `json:"status"`
}

type ConversationRequest struct {
	ConversationId  string `json:"conversationId"`
	ParentMessageId string `json:"parentMessageId"`
	MessageId       string `json:"messageId"`
}

type ChatRequest struct {
	Options       ConversationRequest `json:"options"`
	Prompt        string              `json:"prompt"`
	SystemMessage string              `json:"systemMessage"`
}

type BaseChatMessage struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Detail *BaseChatMessage `json:"detail"`
}

// ChatMessage ConversationResponse
type ChatMessage struct {
	Id              string            `json:"id"`
	Text            string            `json:"text"`
	Role            Role              `json:"role"`
	Name            string            `json:"name"`
	Delta           string            `json:"delta"`
	Detail          ChatMessageDetail `json:"detail"`
	ParentMessageId string            `json:"parentMessageId"`
	ConversationId  string            `json:"conversationId"`
	AudioUrl        string            `json:"audio_url"`
}
type ChatMessageDetail struct {
	Choices interface{} `json:"choices"`
	Created int64       `json:"created"`
	Id      string      `json:"id"`
	Model   string      `json:"model"`
	Object  string      `json:"object"`
	UseAge  interface{} `json:"useage"`
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

type ConfigResp struct {
	TimeoutMs    int32  `json:"timeoutMs"`
	ReverseProxy string `json:"reverseProxy"`
	ApiModel     string `json:"apiModel"`
	SocksProxy   string `json:"socksProxy"`
	HttpsProxy   string `json:"httpsProxy"`
	Valance      string `json:"balance"`
}

type SessionResp struct {
	Auth  bool   `json:"auth"`
	Model string `json:"model"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type Choice struct {
	Delta        Delta       `json:"delta"`
	Index        int32       `json:"index"`
	FinishReason string      `json:"finish_reason"`
	LogProb      interface{} `json:"logprobs"`
	Text         string      `json:"text"`
}

type UseAge struct {
	CompletionTokens int64 `json:"completion_tokens"`
	PromptTokens     int64 `json:"prompt_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type CreateChatCompletionDeltaResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int32    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Delta struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type AsrResponse struct {
	Text string `json:"text"`
}

type AdvanceRequest struct {
	SystemMessage string `json:"systemMessage"`
	Model         string `json:"model"`
}

type AdvanceResponse struct {
	SystemMessage string        `json:"systemMessage"`
	Model         string        `json:"model"`
	Avatar        []Image       `json:"avatar"`
	ModelList     []OptionModel `json:"modelList"`
}

type Image struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

type OptionModel struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type UserInfo struct {
	Avatar      string `json:"avatar"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type History struct {
	Title  string `json:"title"`
	IsEdit bool   `json:"isEdit"`
	Uuid   int    `json:"uuid"`
}

type ChatHistory struct {
	Active       int       `json:"active"`
	UsingContext bool      `json:"usingContext"`
	History      []History `json:"history"`
	Chat         []Chat    `json:"chat"`
}

type Chat struct {
	Uuid int        `json:"uuid"`
	Data []ChatBase `json:"data"`
}

type RequestOptions struct {
	Prompt  string              `json:"prompt"`
	Options ConversationRequest `json:"options"`
}

type ChatBase struct {
	DateTime            string              `json:"dateTime"`
	Text                string              `json:"text"`
	Inversion           bool                `json:"inversion"`
	Error               bool                `json:"error"`
	Loading             bool                `json:"loading"`
	ConversationOptions ConversationRequest `json:"conversation_options"`
	RequestOptions      RequestOptions      `json:"requestOptions"`
	AudioUrl            string              `json:"audio_url"`
}

type Prompt struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PromptListResp struct {
	PromptList []Prompt `json:"promptList"`
}
