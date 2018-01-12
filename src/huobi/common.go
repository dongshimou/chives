package huobi

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"
)

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}
func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
func parseTS2String(ts int64) string {
	return parseTS2Time(ts).Format("2006-01-02 15:04:05")
}
func parseTS2Time(ts int64) time.Time {
	//时间戳 1515408671212 去掉 212
	return time.Unix(ts/1000, 0)
}

func HmacSHA256(message, secret string) string {
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(sig.Sum(nil))
}

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

func createUrl(method, path string, para map[string]interface{}) (string, error) {
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
		para = map[string]interface{}{}
	}
	para[SignatureMethod] = signMethod
	para[SignatureVersion] = signVer
	para[AccessKeyId] = accessKeyID
	para[Timestamp] = timeStamp

	list := []string{}

	for k, v := range para {
		list = append(list, k+"="+url.QueryEscape(fmt.Sprintf("%v", v)))
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
		urls += fmt.Sprintf("%v", v)
	}

	return urls, nil
}
