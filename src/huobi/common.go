package huobi

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
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
	//时间戳 1515408671212 去掉 1212
	return time.Unix(ts/1000, 0)
}

func maxF32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func minF32(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}
