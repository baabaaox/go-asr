package asr

import (
	"flag"
	"testing"
)

var googleAppKey string

func init() {
	flag.StringVar(&googleAppKey, "googleAppKey", "", "Google App Key")
}

func TestGoogleRecognize(t *testing.T) {
	asrEngine := Engine{
		IEngine: &Google{
			APIKey:     googleAppKey,
			Encoding:   1,
			SampleRate: 16000,
			Lang:       "zh",
		},
	}
	text, err := asrEngine.recognize("test.pcm")
	if err != nil {
		t.Error(err)
	}
	if text != "您好" {
		t.Errorf("Text [%v]", text)
	}
}
