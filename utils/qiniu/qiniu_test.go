package qiniu

import "testing"

func TestUploadLocalFile(t *testing.T) {
	UploadLocalFile("test.go2", "./qiniu_test.go")
}
