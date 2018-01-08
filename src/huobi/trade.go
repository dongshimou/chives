package huobi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

func startTrade(c *websocket.Conn, market string) (err error) {
	//初始化订阅
	initTrade(market)
	//发送订阅数据
	c.WriteJSON(trades)
	//获得订阅状态
	err = func() error {
		_, hstr, err := c.ReadMessage()
		if err != nil {
			log.Println("hello ReadMessage err is ", err.Error())
			return err
		}
		hraw, err := GzipDecode(hstr)
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
	}()
	if err != nil {
		return err
	}
	//循环读取订阅
	go func() {
		defer c.Close()
		for {
			//读取 websocket 数据
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("ReadMessage err is ", err.Error())
			}
			// Gzip 解压
			raw, err := GzipDecode(message)
			if err != nil {
				log.Println("Gzip err is ", err.Error())
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
			}
			//展示数据
			TradeShow(&v)
		}
	}()
	return nil
}
func TradeShow(data *Trade) {
	show := func(v *TradeTickData) {
		//时间戳 1515408671212 去掉 1212
		t := time.Unix(v.TS/1000, 0)
		str := t.Format("2006-01-02 15:04:05")
		if v.Direction == BUY {
			str += " : 买入 -->> "
		} else {
			str += " : 卖出 <<-- "
		}
		str += fmt.Sprintf("价格:%f 成交量:%f ", v.Price, v.Amount)
		log.Println(str)
	}
	for _, v := range data.Tick.Data {
		show(&v)
	}
}
