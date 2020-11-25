package asr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Google Cloud Platform
	// https://cloud.google.com/speech-to-text/docs/reference/rest/v1p1beta1/speech/recognize
	Google struct {
		APIKey     string
		Encoding   int32
		SampleRate int32
		Lang       string
	}
)

func (google *Google) prepare() (err error) {
	return
}

func (google *Google) recognize(filename string) (text string, err error) {
	url := "https://speech.googleapis.com/v1p1beta1/speech:recognize?key=" + google.APIKey
	audioData, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	payload := make(map[string]interface{})
	payload["config"] = map[string]interface{}{
		"encoding":        google.Encoding,
		"sampleRateHertz": google.SampleRate,
		"languageCode":    google.Lang,
	}
	payload["audio"] = map[string]string{
		"content": base64.StdEncoding.EncodeToString(audioData),
	}
	payloadMap, _ := json.Marshal(payload)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payloadMap))
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	err = fmt.Errorf("Failed to recognize [%d] %s", response.StatusCode, responseBody)
	if response.StatusCode == 200 {
		type responseStruct struct {
			Results []struct {
				Alternatives []struct {
					Transcript string
				}
			}
		}
		var res responseStruct
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return
		}
		text = res.Results[0].Alternatives[0].Transcript
		err = nil
	}
	return
}
