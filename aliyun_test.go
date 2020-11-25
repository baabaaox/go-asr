package asr

import (
	"flag"
	"testing"
)

var aliyunAppID string
var aliyunAppKey string
var aliyunAppSecret string

func init() {
	flag.StringVar(&aliyunAppID, "aliyunAppID", "", "Aliyun App ID")
	flag.StringVar(&aliyunAppKey, "aliyunAppKey", "", "Aliyun App Key")
	flag.StringVar(&aliyunAppSecret, "aliyunAppSecret", "", "Aliyun App Secret")
}

func TestAliyunRecognize(t *testing.T) {
	asrEngine := Engine{
		IEngine: &Aliyun{
			AppID:           aliyunAppID,
			AccessKeyID:     aliyunAppKey,
			AccessKeySecret: aliyunAppSecret,
			Format:          "pcm",
			SampleRate:      16000,
		},
	}
	err := asrEngine.Prepare()
	if err != nil {
		t.Error(err)
	}
	text, err := asrEngine.Recognize("test.pcm")
	if err != nil {
		t.Error(err)
	}
	if text != "您好" {
		t.Errorf("Text [%v]", text)
	}
}
