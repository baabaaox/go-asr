package asr

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type (
	// Aliyun Cloud Platform
	// https://help.aliyun.com/document_detail/92131.html?spm=a2c4g.11186623.6.576.7b60ac6cXfLxTd
	Aliyun struct {
		AppID           string
		AccessKeyID     string
		AccessKeySecret string
		Format          string
		SampleRate      int32
		token           string
	}
)

func hmacsha1(s, key string) string {
	hashed := hmac.New(sha1.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

func (aliyun *Aliyun) setToken() (err error) {
	host := "http://nls-meta.cn-shanghai.aliyuncs.com/?"
	nonce, _ := uuid.NewRandom()
	params := url.Values{}
	params.Add("AccessKeyId", aliyun.AccessKeyID)
	params.Add("Action", "CreateToken")
	params.Add("Format", "JSON")
	params.Add("RegionId", "cn-shanghai")
	params.Add("SignatureMethod", "HMAC-SHA1")
	params.Add("SignatureNonce", nonce.String())
	params.Add("SignatureVersion", "1.0")
	params.Add("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	params.Add("Version", "2019-02-28")
	queryString := params.Encode()
	stringToSign := "GET&%2F&" + url.QueryEscape(queryString)
	signature := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(hmacsha1(stringToSign, aliyun.AccessKeySecret+"&"))))
	queryStringWithSign := "Signature=" + signature + "&" + queryString
	response, err := http.Get(host + queryStringWithSign)
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	err = fmt.Errorf("Authentication failed [%d] %s", response.StatusCode, responseBody)
	if response.StatusCode == 200 {
		type responseStruct struct {
			Token struct {
				ID string
			}
		}
		var res responseStruct
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return
		}
		aliyun.token = res.Token.ID
		err = nil
	}
	return
}

func (aliyun *Aliyun) prepare() (err error) {
	err = aliyun.setToken()
	return
}

func (aliyun *Aliyun) recognize(filename string) (text string, err error) {
	url := "http://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/asr"
	url = fmt.Sprintf("%s?appkey=%s&format=%s&sample_rate=%d", url, aliyun.AppID, aliyun.Format, aliyun.SampleRate)
	audioData, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(audioData))
	if err != nil {
		return
	}
	request.Header.Add("X-NLS-Token", aliyun.token)
	request.Header.Add("Content-Type", "application/octet-stream")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	err = fmt.Errorf("Failed to recognize [%d] %s", response.StatusCode, responseBody)
	if response.StatusCode == 200 {
		type responseStruct struct {
			Result string
		}
		var res responseStruct
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return
		}
		text = res.Result
		err = nil
	}
	return
}
