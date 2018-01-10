package huobi

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func input() string {

	log.Println("========================================")
	log.Println("输入你要查看的即时数据的交易类型:")
	log.Println("如 : ethusdt ")
	log.Println("查询 eth 与 usdt 的即时行情")
	log.Println("========================================")

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text

}
func initWS() (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.huobi.pro/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil, err
	}
	return c, nil
}
func Run() {
	market := ""
	for {
		c, err := initWS()
		market = input()
		//market = "ethusdt"
		InitDB(market)
		log.Println("所查行情为 :", market)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = startTrade(c, market)
		if err == nil {
			break
		} else {
			log.Println("==========启动失败正在重新启动==========")
			log.Println()
			log.Println()
			log.Println()
		}
	}
	retryTrade := func() {
		c, err := initWS()
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = startTrade(c, market)
		if err != nil {
			log.Println("==========启动失败正在重新启动==========")
		}
	}
	//检查是否关闭
	go func() {
		for {
			if isTradeClose() {
				log.Println("即时行情异常关闭,正在重启")
				go retryTrade()
			}
			time.Sleep(time.Second * 10)
		}
	}()
	for {

		getLegal(CNY_USDT)
		totalmax, totalmin := getTotalMaxMinPrice()
		currmax, currmin := getCurrentMaxMinPrice()
		currAmount := getCurrentAmount()
		setCurrentInit()
		log.Println(fmt.Sprintf("========================当前 %s 的价格======================", market))
		log.Println(fmt.Sprintf("===历史的最高价 : %14.4f ===========历史的最低价 : %14.4f ===", totalmax, totalmin))
		log.Println(fmt.Sprintf("===上分钟最高价 : %14.4f ===========上分钟最低价 : %14.4f ===", currmax, currmin))
		log.Println(fmt.Sprintf("=======================上分钟成交量: %32.4f ===============", currAmount))
		time.Sleep(time.Minute)
	}
}
