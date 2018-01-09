package huobi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync/atomic"
)

var tradeStatus = STATUS_OPEN

const (
	STATUS_CLOSE int32 = 0
	STATUS_OPEN  int32 = 1
)

func isTradeClose() bool {
	return atomic.LoadInt32(&tradeStatus) == STATUS_CLOSE
}
func setTradeClose() {
	atomic.StoreInt32(&tradeStatus, STATUS_CLOSE)
}

func processData(c *websocket.Conn, data chan []byte) {
	for {
		if isTradeClose() {
			return
		}
		zipdata := <-data
		raw, err := GzipDecode(zipdata)
		if err != nil {
			log.Println("hello Gzip err is ", err.Error())
		}
		msg := string(raw)
		//检查是否是 心跳包
		if in := strings.Index(msg, "ping"); in > 0 {
			//log.Println("this is a ping !!!")
			msg = strings.Replace(msg, "ping", "pong", 1)
			c.WriteMessage(websocket.TextMessage, []byte(msg))
			continue
		}
		//获得数据
		v := Trade{}
		err = json.Unmarshal(raw, &v)
		if err != nil {
			log.Println("json parse err is ", err.Error())
		} else {
			//展示数据
			TradeShow(&v)
		}
	}
}
func helloTrade(c *websocket.Conn) (err error) {
	_, zipstr, err := c.ReadMessage()
	if err != nil {
		log.Println("hello ReadMessage err is ", err.Error())
		return err
	}
	hraw, err := GzipDecode(zipstr)
	if err != nil {
		log.Println("hello Gzip err is ", err.Error())
		return err
	}
	h := TradeHello{}
	err = json.Unmarshal(hraw, &h)
	if err != nil {
		log.Println("hello json parse err is ", err.Error())
		return err
	}
	if h.Status != "ok" {
		log.Println("握手失败,请重试")
		log.Println(h.ErrMsg)
		return NewError(0, "握手错误")
	} else {
		log.Println(h.Subbed)
		log.Println("========================交易明细========================")
	}
	return nil
}
func startTrade(c *websocket.Conn, market string) (err error) {
	//初始化订阅
	initTrade(market)
	//发送订阅数据
	log.Println(trades.Sub, "====", trades.ID)
	c.WriteJSON(GetTradeConfig())
	//初始化管道
	rawData := make(chan []byte, 100)
	//获得握手以及订阅状态
	err = helloTrade(c)
	if err != nil {
		c.Close()
		return err
	}
	// 处理数据
	go processData(c, rawData)
	// 循环读取订阅
	go func() {
		defer close(rawData)
		for {
			//读取 websocket 数据
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("ReadMessage err is ", err.Error())
				setTradeClose()
				return
			}
			//放入管道
			rawData <- message
		}
	}()
	return nil
}

func TradeShow(data *Trade) {
	show := func(v *TradeTickData) {
		str := parseTS2String(v.TS)
		if v.Direction == BUY {
			str += " : 买入 -->> "
		} else {
			str += " : 卖出 <<-- "
		}
		str += fmt.Sprintf("价格:%f 成交量:%f ", v.Price, v.Amount)
		log.Println(str)
	}
	log.Println(fmt.Sprintf("====%s================ 交易ID : %d ====",
		parseTS2String(data.Tick.TS), data.Tick.ID))
	for _, v := range data.Tick.Data {
		show(&v)
	}
}
