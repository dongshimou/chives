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

type MarketDetail struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}
type MarketKline struct {
	Req string `json:"req"`
	ID  string `json:"id"`
}

var marketDetail = MarketDetail{}
var marketKline = MarketKline{}

var tradeTemplate = MarketDetail{
	Sub: "market.symbol.trade.detail",
	ID:  "id",
}
var klineTemplate = MarketKline{
	Req: "market.symbol.kline.1min",
	ID:  "id",
}

func initTrade(s string) {
	marketDetail = tradeTemplate
	marketDetail.Sub = strings.Replace(marketDetail.Sub, "symbol", s, 1)
	rand.Seed(time.Now().UnixNano())
	x := rand.Uint64()
	marketDetail.ID = strings.Replace(marketDetail.ID, "id", fmt.Sprintf("id%d", x), 1)
}
func initMarketKline(s string) {
	marketKline = klineTemplate
	marketKline.Req = strings.Replace(marketKline.Req, "symbol", s, 1)
	rand.Seed(time.Now().UnixNano())
	x := rand.Uint64()
	marketKline.ID = strings.Replace(marketKline.ID, "id", fmt.Sprintf("id%d", x), 1)
}
func GetMarketDetailConfig() *MarketDetail {
	return &marketDetail
}
func getMarketKlineConfig() *MarketKline {
	return &marketKline
}

type KeyConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}
type DBConfig struct {
	DBSave     bool   `json:"db_save"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBDatabase string `json:"db_database"`
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
