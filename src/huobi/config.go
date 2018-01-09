package huobi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
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

func readConfig(fileName string, configObj interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	content, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(content, configObj)
	return err
}
