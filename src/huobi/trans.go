package huobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"time"
)

const (
	NEXT_LINE = "\n"

	GET  = "GET"
	POST = "POST"
)

//GET, "/v1/order/orders", nil
func createUrl(method, path string, para []string) (string, error) {
	keys := Key{}
	err := readConfig("./key", &keys)
	if err != nil {
		return "", err
	}
	access := keys.AccessKey

	scress := keys.SecretKey

	host := "api.huobi.pro"

	req := ""

	req += method
	req += NEXT_LINE
	req += host
	req += NEXT_LINE
	req += path
	req += NEXT_LINE

	//parameter
	signMethod := "SignatureMethod=HmacSHA256"
	signVer := "SignatureVersion=2"
	accessKeyID := "AccessKeyId=" + access
	timeStamp := fmt.Sprintf("Timestamp=%04d-%02d-%02dT%02d%%3A%02d%%3A%02d",
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	parameter := ""
	para = append(para, signMethod)
	para = append(para, signVer)
	para = append(para, accessKeyID)
	para = append(para, timeStamp)
	sort.Strings(para)
	for i, v := range para {
		if i != 0 {
			parameter += "&"
		}
		parameter += v
	}

	unsignmsg := req + parameter
	//hmacsha256加密
	sig := hmac.New(sha256.New, []byte(scress))
	sig.Write([]byte(unsignmsg))

	base64str := base64.StdEncoding.EncodeToString(sig.Sum(nil))

	url := "https://" + host + path + parameter + "&Signature=" + base64str

	return url, nil
}
