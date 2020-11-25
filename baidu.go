package asr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type (
	// Baidu Cloud Platform
	// https://ai.baidu.com/ai-doc/SPEECH/6k38lxp0r
	Baidu struct {
		ClientID     string
		ClientSecret string
		Format       string
		token        string
	}
)

func (baidu *Baidu) setToken() (err error) {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	url = fmt.Sprintf("%s?grant_type=client_credentials&client_id=%s&client_secret=%s", url, baidu.ClientID, baidu.ClientSecret)
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	err = fmt.Errorf("Authentication failed [%d] %s", response.StatusCode, responseBody)
	if response.StatusCode == 200 {
		type responseStruct struct {
			AccessToken string `json:"access_token"`
		}
		var res responseStruct
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return
		}
		baidu.token = res.AccessToken
		err = nil
	}
	return
}

func (baidu *Baidu) prepare() (err error) {
	err = baidu.setToken()
	return
}

func (baidu *Baidu) recognize(filename string) (text string, err error) {
	url := "https://vop.baidu.com/pro_api"
	audioData, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	cuid, _ := uuid.NewRandom()
	payload := make(map[string]interface{})
	payload["format"] = baidu.Format
	payload["rate"] = 16000
	payload["channel"] = 1
	payload["cuid"] = cuid
	payload["token"] = baidu.token
	payload["dev_pid"] = 80001
	payload["len"] = len(audioData)
	payload["speech"] = base64.StdEncoding.EncodeToString(audioData)
	payloadMap, _ := json.Marshal(payload)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payloadMap))
	if err != nil {
		return
	}
	responseBody, _ := ioutil.ReadAll(response.Body)
	err = fmt.Errorf("Failed to recognize [%d] %s", response.StatusCode, responseBody)
	if response.StatusCode == 200 {
		type responseStruct struct {
			Result []string
		}
		var res responseStruct
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return
		}
		text = res.Result[0]
		err = nil
	}
	return
}
