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

	totalMaxPrice float64 = MAX_FLOAT_INIT
	totalMinPrice float64 = MIN_FLOAT_INIT
	//totalAve      float64 = 0
	//totalCount    int64   = 0

	currMaxPrice   float64 = MAX_FLOAT_INIT
	currMinPrice   float64 = MIN_FLOAT_INIT
	currAve        float64 = 0
	currCount      float64 = 0
	currentAmount  float64 = 0
	currAmountBuy  float64 = 0
	currAmountSell float64 = 0

	currUpCount   float64 = 0
	currDownCount float64 = 0

	lastPrice float64 = 0
)

const (
	STATUS_CLOSE int32 = 0
	STATUS_OPEN  int32 = 1

	MIN_FLOAT_INIT = math.MaxFloat64
	MAX_FLOAT_INIT = -math.MaxFloat64
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
	currMaxPrice = math.MinInt32
	currMinPrice = math.MaxInt32
	currentAmount, currAmountBuy, currAmountSell = 0, 0, 0
	currAve = 0
	currCount, currUpCount, currDownCount = 0, 0, 0
}
func getTotalMaxMinPrice() (max, min float64) {
	return totalMaxPrice, totalMinPrice
}
func getCurrentMaxMinPrice() (max, min float64) {
	return currMaxPrice, currMinPrice
}
func getCurrentAmount() float64 {
	return currentAmount
}
func getAveragePrice() float64 {
	return currAve
}
func getLastPrice() float64 {
	return lastPrice
}
func getAverageCount() float64 {
	return currCount
}
func getUpDownCount() (float64, float64) {
	return currUpCount, currDownCount
}
func isUpOrDown() int {
	temp := math.Abs(currUpCount - currDownCount)
	if temp < 0 {
		return -1
	} else if temp/currCount < 0.1 {
		return 0
	} else {
		return 1
	}
}

func pong(c *websocket.Conn, msg string) {
	msg = strings.Replace(msg, "ping", "pong", 1)
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}
func hello(raw []byte) {
	h := TradeHello{}
	err := json.Unmarshal(raw, &h)
	if err != nil {
		log.Println("hello json parse err is ", err.Error())
	}
	if h.Status != "ok" {
		log.Println("握手失败,请重试")
		log.Println(h.ErrMsg)
	} else {
		log.Println(h.Subbed)
		log.Println("========================交易明细========================")
	}
}
func process(raw []byte) *TradeTick {
	//获得数据
	v := Trade{}
	//log.Println(string(raw))
	err := json.Unmarshal(raw, &v)
	if err != nil {
		log.Println("json parse err is ", err.Error())
	} else {
		//展示数据

		TradeShow(&v)
		if IsDBSave() {
			return &v.Tick
		}
	}
	return nil
}
func recv(c *websocket.Conn, data chan []byte) {
	tick := make(chan TradeTick, 1024)
	defer close(tick)
	if IsDBSave() {
		go saveTick(tick)
	}
	for {
		zipdata, ok := <-data
		if !ok {
			log.Println("recv 管道 不存在数据")
			return
		}
		raw, err := GzipDecode(zipdata)
		if err != nil {
			log.Println("hello Gzip err is ", err.Error())
		}
		msg := string(raw)

		if in := strings.Index(msg, "ping"); in > 0 {
			pong(c, msg)
			continue
		} else if in := strings.Index(msg, "status"); in > 0 {
			hello(raw)
		} else {
			t := process(raw)
			if t != nil {
				tick <- *t
			}
		}

	}
}

func helloTrade(c *websocket.Conn) (err error) {
	for {
		_, zipstr, err := c.ReadMessage()
		if err != nil {
			log.Println("hello ReadMessage err is ", err.Error())
			return err
		}
		hraw, err := GzipDecode(zipstr)
		log.Println(string(hraw))
		if err != nil {
			log.Println("hello Gzip err is ", err.Error())
			return err
		}
		msg := string(hraw)
		if in := strings.Index(msg, "ping"); in > 0 {
			pong(c, msg)
			continue
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
}
func startTrade(c *websocket.Conn, market string) (err error) {
	//初始化订阅
	initTrade(market)
	//设置开始
	setTradeOpen()
	//发送订阅数据
	log.Println(trades.Sub, "====", trades.ID)
	c.WriteJSON(GetTradeConfig())

	//初始化管道
	rawData := make(chan []byte, 1024)
	// 处理数据
	go recv(c, rawData)
	// 循环读取订阅
	go func() {
		defer setTradeClose()
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

		currCount++
		temp := (currCount - 1) / currCount
		currAve = currAve*temp + v.Price/currCount

		if lastPrice < v.Price {
			//统计上涨
			currUpCount++
		} else if lastPrice > v.Price {
			//统计下跌
			currDownCount++
		}
		lastPrice = v.Price

		currMinPrice = math.Min(v.Price, currMinPrice)
		currMaxPrice = math.Max(v.Price, currMaxPrice)
		totalMinPrice = math.Min(totalMinPrice, currMinPrice)
		totalMaxPrice = math.Max(totalMaxPrice, currMaxPrice)
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
