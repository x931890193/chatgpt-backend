package xunfei

import (
	"chatgpt-backend/logger"
	"io"
	"os"
	"testing"
)

func TestAsrRealtime(t *testing.T) {
	audio, _ := os.Open("./111.wav")
	content, err := io.ReadAll(audio)
	if err != nil {
		return
	}
	res := AsrRealtime(content)
	logger.Info.Println("res", res)
}

func TestTTS(t *testing.T) {
	res := TtsClient.TTS("你好你叫什么名字吗")
	println(res)
	f, err := os.Create("111.mp3")
	if err != nil {
		return
	}
	f.Write(res)
}

func TestAsrStream(t *testing.T) {
	audio, _ := os.Open("16k_10.pcm")
	//content, _ := io.ReadAll(audio)
	res, err := AsrStreamClient.Asr(audio)
	if err != nil {
		println(err.Error())
	}
	println(111, res)
}
