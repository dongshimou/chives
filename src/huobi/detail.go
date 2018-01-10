package huobi

import (
	"fmt"
	"log"
	"math"
	"sync/atomic"
)

const (
	STATUS_CLOSE int32 = 0
	STATUS_OPEN  int32 = 1
)

var (
	tradeStatus = STATUS_OPEN

	totalMaxPrice float64 = MAX_FLOAT_INIT
	totalMinPrice float64 = MIN_FLOAT_INIT

	currMaxPrice   float64 = MAX_FLOAT_INIT
	currMinPrice   float64 = MIN_FLOAT_INIT
	currAve        float64 = 0
	currCount      float64 = 0
	currentAmount  float64 = 0
	currAmountBuy  float64 = 0
	currAmountSell float64 = 0

	currUpCount      float64 = 0
	currBalanceCount float64 = 0
	currDownCount    float64 = 0

	lastPrice float64 = 0
)

func isDetailClose() bool {
	return atomic.LoadInt32(&tradeStatus) == STATUS_CLOSE
}
func setDetailClose() {
	atomic.StoreInt32(&tradeStatus, STATUS_CLOSE)
}
func setDetailOpen() {
	atomic.StoreInt32(&tradeStatus, STATUS_OPEN)
}
func setCurrentInit() {
	currMaxPrice = math.MinInt32
	currMinPrice = math.MaxInt32
	currentAmount, currAmountBuy, currAmountSell = 0, 0, 0
	currAve = 0
	currCount, currUpCount, currBalanceCount, currDownCount = 0, 0, 0, 0
}
func getTSMinuteTS(ts int64) int64 {
	return ts / 1000 / 60 * 60
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
func getUpDownCount() (up, ba, down float64) {
	return currUpCount, currBalanceCount, currDownCount
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
func showDetail(data *Trade) {
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
		} else {
			//震荡
			currBalanceCount++
		}
		lastPrice = v.Price

		currMinPrice = math.Min(v.Price, currMinPrice)
		currMaxPrice = math.Max(v.Price, currMaxPrice)
		totalMinPrice = math.Min(totalMinPrice, currMinPrice)
		totalMaxPrice = math.Max(totalMaxPrice, currMaxPrice)
		currentAmount += v.Amount

		str += fmt.Sprintf("价格:%15.4f 成交量:%13.4f ", v.Price, v.Amount)
		log.Println(str)
		{
			//kline 数据

			//可以通过api请求
		}
	}
	log.Println(fmt.Sprintf("====%s================ 交易ID : %d ====",
		parseTS2String(data.Tick.TS), data.Tick.ID))

	for _, v := range data.Tick.Data {
		show(&v)
	}
}
