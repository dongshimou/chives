package huobi

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type TradeDetail struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

var trades = TradeDetail{}

var tradeTemplate = TradeDetail{
	Sub: "market.symbol.trade.detail",
	ID:  "id",
}

func initTrade(s string) {
	trades = tradeTemplate
	trades.Sub = strings.Replace(trades.Sub, "symbol", s, 1)
	rand.Seed(time.Now().UnixNano())
	x := rand.Uint64()
	trades.ID = strings.Replace(trades.ID, "id", fmt.Sprintf("id%d", x), 1)
}
func GetTradeConfig() *TradeDetail {
	return &trades
}
