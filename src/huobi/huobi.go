package huobi

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"strconv"
)

var (
	Markets = map[string]bool{}
)

func addMarket(s string) {
	Markets[s] = true
}
func delMarket(s string) {
	Markets[s] = false
}
func isMarket(s string) bool {
	v, exist := Markets[s]
	if !exist || !v {
		return false
	}
	return true
}

func init() {
	initConfig()
	getAccountID()
}

func goin() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

func Run() {
	for {
		log.Println("input : =========================")
		log.Println(" show_cny     :")
		log.Println(" create_order :")
		log.Println(" open_detail  :")
		log.Println(" close_show   :")
		log.Println(" open_show    :")
		log.Println("        =========================")

		cmd := goin()
		switch cmd {
		case "show_cny":
			 startTools_GetCNY(initToolsGetCNY())
		case "create_order":
			 startTools_PostOrder(initTools_PostOrder())
		case "open_detail":
			 startTools_Detail(initTools_Detail())
			}
	}
}
func initToolsGetCNY() (int, uint64) {
	log.Println("input market : usdt , btc . It will show exchange price with CNY ")
	in := goin()
	market := 2
	if in == "usdt" {
		market = CNY_USDT
	} else if in == "btc" {
		market = CNY_BTC
	}
	log.Println("input a number (x). It will show the market every (x) second ")
	in = goin()
	x, err := strconv.ParseUint(in, 10, 32)
	if err != nil {
		x = 60
	}
	return market, x
}
func startTools_GetCNY(mar int, dur uint64) {
	for {
		getLegal(mar)
		time.Sleep(time.Second * time.Duration(dur))
	}
}
func initTools_PostOrder() (string, string, float64, float64) {
	for {
		log.Println("\n\n\n\n\n==============================")
		log.Println("====输入交易对==== 如 : ethusdt")
		market := goin()
		log.Println("===输入交易方式=== 如 : =======")
		log.Println("    bm        市价买(buy-market) ")
		log.Println("    sm        市价卖(sell-market)")
		log.Println("    bl        限价买(buy-limit)  ")
		log.Println("    sl        限价卖(sell-limit) ")
		log.Println("==============================")
		ordertype := goin()
		amount, price := "0", "1"
		switch ordertype {
		case "bm":
			ordertype = BUY_MARKET
			log.Println("=======输入 买多少钱的=======")
			amount = goin()
			log.Println(fmt.Sprintf("%s     市价买入     花费:%s ", market, amount))
			log.Println("=========================================================")
		case "sm":
			ordertype = SEL_MARKET
			log.Println("=======输入 卖多少数量=======")
			amount = goin()
			log.Println(fmt.Sprintf("%s     市价卖出     数量:%s ", market, amount))
			log.Println("=========================================================")
		case "bl":
			ordertype = BUY_LIMIT
			log.Println("=======输入 买入数量=========")
			amount = goin()
			log.Println("=======输入 买入单价=========")
			price = goin()
			log.Println(fmt.Sprintf("%s     限价买入     数量:%s 单价:%s ", market, amount, price))
			log.Println("=========================================================")
		case "sl":
			ordertype = SELL_LIMIT
			log.Println("=======输入 卖出数量=========")
			amount = goin()
			log.Println("=======输入 卖出单价=========")
			price = goin()
			log.Println(fmt.Sprintf("%s     限价卖出     数量:%s 单价:%s ", market, amount, price))
			log.Println("=========================================================")
		}
		a, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Println("amount error")
			continue
		}
		p, err := strconv.ParseFloat(price, 64)
		if err != nil {
			log.Println("price error")
			continue
		}
		return market, ordertype, a, p
	}
}
func startTools_PostOrder(market, ordertype string, a, p float64) {
	createOrder(market, ordertype, a, p)
}
func initTools_Detail() string {
	setDetailClose()
	log.Println("========================================")
	log.Println("输入你要查看的即时数据的交易类型:")
	log.Println("如 : ethusdt ")
	log.Println("查询 eth 与 usdt 的即时行情")
	log.Println("========================================")
	return goin()
}
func startTools_Detail(market string) {

	retryDetail := func() {

		c, err := initWS()
		if err != nil {
			log.Println("初始化 websocket 错误")
			return
		}
		err=InitDB(market)
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = startWS(c, market)
		if err != nil {
			log.Println("==========启动失败正在重新启动==========\n\n\n")
		}
		addMarket(market)
	}

	//检查是否关闭
	go func() {
		for {
			if isDetailClose() {
				log.Println("正在启动 即时行情")
				go retryDetail()
			}
			time.Sleep(time.Second * 3)
		}
	}()

	//myideal(market)

	for {
		time.Sleep(time.Minute)
	}
}
