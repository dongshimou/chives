package huobi

import (
	"bufio"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
	"time"
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
	for {
		if isTradeClose() {
			log.Println("即时行情异常关闭,正在重启")
			go retryTrade()
		}
		GetLegal(CNY_USDT)
		time.Sleep(time.Second * 10)
	}
}
