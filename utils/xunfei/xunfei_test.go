package xunfei

import (
	"chatgpt-backend/logger"
	"io"
	"os"
	"testing"
)

func TestAsr(t *testing.T) {
	audio, _ := os.Open("./111.wav")
	content, err := io.ReadAll(audio)
	if err != nil {
		return
	}
	res := Asr(content)
	logger.Info.Println("res", res)
}

func TestTTS(t *testing.T) {
	res := TtsClient.TTS("你好你叫什么名字吗")
	println(res)
	f, err := os.Create("111.wav")
	if err != nil {
		return
	}
	f.Write(res)
}
