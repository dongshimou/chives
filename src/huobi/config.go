package huobi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"log"
)

type addMarketDetail struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}
type delMarketDetail struct {
	Unsub string `json:"unsub"`
	ID    string `json:"id"`
}
type addMarketKline struct {
	Req string `json:"req"`
	ID  string `json:"id"`
}

var addDetailTemplate = addMarketDetail{
	Sub: "market.symbol.trade.detail",
	ID:  "id",
}
var delDetailTemplate = delMarketDetail{
	Unsub: "market.symbol.trade.detail",
	ID:    "id",
}
var addKlineTemplate = addMarketKline{
	Req: "market.symbol.kline.1min",
	ID:  "id",
}

func getRandomInt64() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Uint64())
}
func generateAddMarketDetail(s string) addMarketDetail {
	marketDetail := addDetailTemplate
	marketDetail.Sub = strings.Replace(marketDetail.Sub, "symbol", s, 1)
	marketDetail.ID = strings.Replace(marketDetail.ID, "id", fmt.Sprintf("id%d", getRandomInt64), 1)
	return marketDetail
}
func generateDelMarketDetail(s string) delMarketDetail {
	marketDetail := delDetailTemplate
	marketDetail.Unsub = strings.Replace(marketDetail.Unsub, "symbol", s, 1)
	marketDetail.ID = strings.Replace(marketDetail.ID, "id", fmt.Sprintf("id%d", getRandomInt64), 1)
	return marketDetail
}
func generateAddMarketKline(s string) addMarketKline {
	marketKline := addKlineTemplate
	marketKline.Req = strings.Replace(marketKline.Req, "symbol", s, 1)
	marketKline.ID = strings.Replace(marketKline.ID, "id", fmt.Sprintf("id%d", getRandomInt64), 1)
	return marketKline
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
type Config struct {
	Database DBConfig `json:"database"`
	TransKey KeyConfig `json:"trans_key"`
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
func getDBConfig()*DBConfig{
	return &chivesConfig.Database
}
func getKeyConfig()*KeyConfig{
	return &chivesConfig.TransKey
}
var(
	chivesConfig = Config{}
)
func initConfig(){
	err:=readConfig("./config.json",&chivesConfig)
	if err!=nil{
		log.Println(err.Error())
	}
}