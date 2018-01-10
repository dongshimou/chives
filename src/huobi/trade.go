package huobi

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var (
	tradeStatus = STATUS_OPEN

	MIN_FLOAT_INIT = math.MaxFloat64
	MAX_FLOAT_INIT = -math.MaxFloat64

	maxPrice float64 = MAX_FLOAT_INIT
	minPrice float64 = MIN_FLOAT_INIT

	currentMax float64 = MAX_FLOAT_INIT
	currentMin float64 = MIN_FLOAT_INIT

	currentAmount  float64 = 0
	currAmountBuy  float64 = 0
	currAmountSell float64 = 0
)

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
func setTradeOpen() {
	atomic.StoreInt32(&tradeStatus, STATUS_OPEN)
}
func setCurrentInit() {
	currentMax = math.MinInt32
	currentMin = math.MaxInt32
	currentAmount, currAmountBuy, currAmountSell = 0, 0, 0
}
func getTotalMaxMinPrice() (max, min float64) {
	return maxPrice, minPrice
}
func getCurrentMaxMinPrice() (max, min float64) {
	return currentMax, currentMin
}
func getCurrentAmount() float64 {
	return currentAmount
}
func processData(c *websocket.Conn, data chan []byte) {
	tick := make(chan TradeTick, 1024)
	defer close(tick)
	if IsDBSave() {
		go saveTick(tick)
	}
	for {
		zipdata, ok := <-data
		if !ok {
			return
		}
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
		//log.Println(string(raw))
		err = json.Unmarshal(raw, &v)
		if err != nil {
			log.Println("json parse err is ", err.Error())
		} else {
			//展示数据
			if IsDBSave() {
				tick <- v.Tick
			}
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
		log.Println(string(hraw))
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
	//设置开始
	setTradeOpen()
	//发送订阅数据
	log.Println(trades.Sub, "====", trades.ID)
	c.WriteJSON(GetTradeConfig())
	//获得握手以及订阅状态
	err = helloTrade(c)
	if err != nil {
		c.Close()
		return err
	}
	//初始化管道
	rawData := make(chan []byte, 1024)
	// 处理数据
	go processData(c, rawData)
	// 循环读取订阅
	go func() {
		defer setTradeClose()
		defer close(rawData)
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

func TradeShow(data *Trade) {
	show := func(v *TradeTickData) {
		str := parseTS2String(v.TS)
		if v.Direction == BUY {
			str += " : 买入 -->> "
			currAmountBuy += v.Amount
		} else {
			str += " : 卖出 <<-- "
			currAmountSell += v.Amount
		}
		currentMin = math.Min(v.Price, currentMin)
		currentMax = math.Max(v.Price, currentMax)
		minPrice = math.Min(minPrice, currentMin)
		maxPrice = math.Max(maxPrice, currentMax)
		currentAmount += v.Amount
		str += fmt.Sprintf("价格:%15.4f 成交量:%13.4f ", v.Price, v.Amount)
		log.Println(str)
	}
	log.Println(fmt.Sprintf("====%s================ 交易ID : %d ====",
		parseTS2String(data.Tick.TS), data.Tick.ID))

	for _, v := range data.Tick.Data {
		show(&v)
	}
}
