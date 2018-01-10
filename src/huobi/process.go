package huobi

import (
	"encoding/json"
	"log"
	"math"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	MIN_FLOAT_INIT = math.MaxFloat64
	MAX_FLOAT_INIT = -math.MaxFloat64
)

func pong(c *websocket.Conn, msg string) {
	msg = strings.Replace(msg, "ping", "pong", 1)
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}
func hello(raw []byte) {
	h := HelloWS{}
	err := json.Unmarshal(raw, &h)
	if err != nil {
		log.Println("helloDetail json parse err is ", err.Error())
	}
	if h.Status != "ok" {
		log.Println("握手失败,请重试")
		log.Println(h.ErrMsg)
	} else {
		log.Println(h.Subbed)
		log.Println("========================交易明细========================")
	}
}
func processDetail(raw []byte) *TradeTick {
	//获得数据
	v := Trade{}
	//log.Println(string(raw))
	err := json.Unmarshal(raw, &v)
	if err != nil {
		log.Println("json parse err is ", err.Error())
	} else {
		//展示数据

		showDetail(&v)
		if IsDBSave() {
			return &v.Tick
		}
	}
	return nil
}
func processKline(raw []byte) *Kline {
	v := Kline{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		log.Println("json parse err is ", err.Error())
	} else {
		//展示数据

		if IsDBSave() {
			return &v
		}
	}
	return nil
}
func processRaw(c *websocket.Conn, data chan []byte) {
	tick := make(chan TradeTick, 1024)
	defer close(tick)
	if IsDBSave() {
		go saveTick(tick)
	}
	for {
		zipdata, ok := <-data
		if !ok {
			log.Println("processRaw 管道 不存在数据")
			return
		}
		raw, err := GzipDecode(zipdata)
		if err != nil {
			log.Println("processRaw Gzip err is ", err.Error())
		}
		msg := string(raw)
		//log.Println(msg)

		findStr := func(str string) bool {
			in := strings.Index(msg, str)
			return in > 0
		}
		if findStr("ping") {
			pong(c, msg)
			continue
		} else if findStr("status") {
			hello(raw)
		} else {
			if findStr("detail") {
				t := processDetail(raw)
				if t != nil {
					tick <- *t
				}
			} else if findStr("kline") {
				t := processKline(raw)
				if t != nil {

				}
			}
		}

	}
}
func sendDetail(c *websocket.Conn, market string) {
	//初始化订阅
	initTrade(market)
	//设置开始
	setDetailOpen()
	//发送订阅数据
	log.Println(trades.Sub, "====", trades.ID)
	c.WriteJSON(GetTradeConfig())
}
func startWS(c *websocket.Conn, market string) (err error) {

	sendDetail(c, market)
	//sendKline(c,market)

	//初始化管道
	rawData := make(chan []byte, 1024)
	// 处理数据
	go processRaw(c, rawData)
	// 循环读取订阅
	go func() {
		defer setDetailClose()
		defer close(rawData)
		defer c.Close()
		for {
			//读取 websocket 数据
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("ReadMessage err is ", err.Error())
				return
			}
			//放入管道
			rawData <- message
		}
	}()

	return nil
}
