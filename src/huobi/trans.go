package huobi

import (
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	NEXT_LINE = "\n"

	GET  = "GET"
	POST = "POST"

	SignatureMethod  = "SignatureMethod"
	SignatureVersion = "SignatureVersion"
	AccessKeyId      = "AccessKeyId"
	Timestamp        = "Timestamp"
	Signature        = "Signature"
)

//GET, "/v1/order/orders", nil
func createUrl(method, path string, para map[string]string) (string, error) {
	keys := KeyConfig{}
	err := readConfig("./transkey", &keys)
	if err != nil {
		return "", err
	}
	access := keys.AccessKey

	secret := keys.SecretKey

	host := "api.huobi.pro"

	req := ""

	req += method
	req += NEXT_LINE
	req += host
	req += NEXT_LINE
	req += path
	req += NEXT_LINE

	//parameter
	signMethod := "HmacSHA256"
	signVer := "2"
	accessKeyID := access

	timeStamp := time.Now().UTC().Format(time.RFC3339)
	timeStamp = strings.TrimRight(timeStamp, "Z")
	if para == nil {
		para = map[string]string{}
	}
	para[SignatureMethod] = signMethod
	para[SignatureVersion] = signVer
	para[AccessKeyId] = accessKeyID
	para[Timestamp] = timeStamp

	list := []string{}

	for k, v := range para {
		list = append(list, k+"="+url.QueryEscape(v))
	}
	sort.Strings(list)

	parameter := ""
	for i, v := range list {
		if i != 0 {
			parameter += "&"
		}
		parameter += v
	}

	unsignmsg := req + parameter
	//hmacsha256加密

	base64str := HmacSHA256(unsignmsg, secret)
	base64str = url.QueryEscape(base64str)
	para[Signature] = base64str

	urls := "https://" + host + path + "?"

	first := true
	for k, v := range para {
		if !first {
			urls += "&"
		}
		first = false
		urls += k
		urls += "="
		urls += v
	}

	return urls, nil
}
