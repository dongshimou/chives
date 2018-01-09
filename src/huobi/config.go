package huobi

import (
	"strings"
)

type TradeDetail struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

var trades = TradeDetail{}

var tradeTemplate = TradeDetail{
	Sub: "market.symbol.trade.detail",
	ID:  "id233333333333322",
}

func initTrade(s string) {
	trades = tradeTemplate
	trades.Sub = strings.Replace(trades.Sub, "symbol", s, 1)
}
func GetTradeConfig() *TradeDetail {
	return &trades
}
