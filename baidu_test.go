package asr

import (
	"flag"
	"testing"
)

var baiduAppKey string
var baiduAppSecret string

func init() {
	flag.StringVar(&baiduAppKey, "baiduAppKey", "", "Baidu App Key")
	flag.StringVar(&baiduAppSecret, "baiduAppSecret", "", "Baidu App Secret")
}

func TestBaiduRecognize(t *testing.T) {
	asrEngine := Engine{
		IEngine: &Baidu{
			ClientID:     baiduAppKey,
			ClientSecret: baiduAppSecret,
			Format:       "wav",
		},
	}
	err := asrEngine.prepare()
	if err != nil {
		t.Error(err)
	}
	text, err := asrEngine.recognize("test.wav")
	if err != nil {
		t.Error(err)
	}
	if text != "您好。" {
		t.Errorf("Text [%v]", text)
	}
}
