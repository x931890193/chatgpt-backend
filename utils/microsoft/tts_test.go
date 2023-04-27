package microsoft

import (
	"chatgpt-backend/logger"
	"io/ioutil"
	"testing"
)

func TestGetVoiceList(t *testing.T) {
	GetVoiceList()
}

func TestGetToken(t *testing.T) {
	GetToken()
}

func TestTTs(t *testing.T) {
	tts, err := TTS("你好吗， 请问你叫什么")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	err = ioutil.WriteFile("output.wav", tts, 0644)
	println(tts)
}
