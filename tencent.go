package asr

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type (
	// Tencent Cloud Platform
	// https://cloud.tencent.com/document/product/1093/37308
	Tencent struct {
		SecretID  string
		SecretKey string
		Format    string
		Lang      string
	}
)

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

// Prepare ...
func (tencent *Tencent) Prepare() (err error) {
	return
}

// Recognize ...
func (tencent Tencent) Recognize(filename string) (text string, err error) {
	host := "asr.tencentcloudapi.com"
	algorithm := "TC3-HMAC-SHA256"
	service := "asr"
	action := "SentenceRecognition"
	version := "2019-06-14"
	timestamp := time.Now().Unix()
	usrAudioKey, _ := uuid.NewRandom()
	audioData, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	// prepare payload
	payload := make(map[string]interface{})
	payload["ProjectId"] = 0
	payload["SubServiceType"] = 2
	payload["EngSerViceType"] = tencent.Lang
	payload["SourceType"] = 1
	payload["VoiceFormat"] = tencent.Format
	payload["UsrAudioKey"] = usrAudioKey
	payload["Data"] = base64.StdEncoding.EncodeToString(audioData)
	payload["DataLen"] = len(audioData)
	payload["FilterPunc"] = 2
	payloadMap, _ := json.Marshal(payload)
	// step 1: build canonical request string
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json; charset=utf-8\n" + "host:" + host + "\n"
	signedHeaders := "content-type;host"
	hashedRequestPayload := sha256hex(string(payloadMap))
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", httpRequestMethod, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, hashedRequestPayload)
	// step 2: build string to sign
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
	hashedCanonicalRequest := sha256hex(canonicalRequest)
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s", algorithm, timestamp, credentialScope, hashedCanonicalRequest)
	// step 3: sign string
	secretDate := hmacsha256(date, "TC3"+tencent.SecretKey)
	secretService := hmacsha256(service, secretDate)
	secretSigning := hmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(hmacsha256(string2sign, secretSigning)))
	// step 4: build authorization
	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", algorithm, tencent.SecretID, credentialScope, signedHeaders, signature)
	// step 5: request
	client := &http.Client{}
	request, err := http.NewRequest("POST", fmt.Sprintf("https://%s", host), bytes.NewBuffer(payloadMap))
	if err != nil {
		return
	}
	request.Header.Set("Authorization", authorization)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Host", host)
	request.Header.Set("X-TC-Action", action)
	request.Header.Set("X-TC-Version", version)
	request.Header.Set("X-TC-Timestamp", fmt.Sprintf("%v", timestamp))
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	err = fmt.Errorf("Failed to recognize [%d] %s", response.StatusCode, responseBody)
	if response.StatusCode == 200 {
		type responseStruct struct {
			Response struct {
				Result string
			}
		}
		var res responseStruct
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return
		}
		text = res.Response.Result
		err = nil
	}
	return
}
