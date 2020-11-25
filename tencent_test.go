package asr

import (
	"flag"
	"testing"
)

var tencentAppKey string
var tencentAppSecret string

func init() {
	flag.StringVar(&tencentAppKey, "tencentAppKey", "", "Tencent App Key")
	flag.StringVar(&tencentAppSecret, "tencentAppSecret", "", "Tencent App Secret")
}

func TestTencentRecognize(t *testing.T) {
	asrEngine := Engine{
		IEngine: &Tencent{
			SecretID:  tencentAppKey,
			SecretKey: tencentAppSecret,
			Format:    "wav",
			Lang:      "16k_zh",
		},
	}
	text, err := asrEngine.recognize("test.wav")
	if err != nil {
		t.Error(err)
	}
	if text != "您好" {
		t.Errorf("Text [%v]", text)
	}
}
