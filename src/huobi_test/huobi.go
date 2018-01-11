package huobi_test

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"
)

func Run() {

	getAccountID()

	for {
		log.Println("\n\n\n\n\n==============================")
		log.Println("====输入交易对==== 如 : ethusdt")
		market := input()
		log.Println("===输入交易方式=== 如 : =======")
		log.Println("    bm        市价买(buy-market) ")
		log.Println("    sm        市价卖(sell-market)")
		log.Println("    bl        限价买(buy-limit)  ")
		log.Println("    sl        限价卖(sell-limit) ")
		log.Println("==============================")
		ordertype := input()
		amount,price:="0","1"
		ordertype=strings.Trim(ordertype," ")
		ordertype=strings.ToLower(ordertype)
		switch ordertype {
		case "bm":
			ordertype=BUY_MARKET
			log.Println("=======输入 买多少钱的=======")
			amount = input()
			log.Println(fmt.Sprintf("%s     市价买入     花费:%s ",market,amount))
			log.Println("=========================================================")
		case "sm":
			ordertype=SEL_MARKET
			log.Println("=======输入 卖多少数量=======")
			amount = input()
			log.Println(fmt.Sprintf("%s     市价卖出     数量:%s ",market,amount))
			log.Println("=========================================================")
		case "bl":
			ordertype=BUY_LIMIT
			log.Println("=======输入 买入数量=========")
			amount = input()
			log.Println("=======输入 买入单价=========")
			price = input()
			log.Println(fmt.Sprintf("%s     限价买入     数量:%s 单价:%s ",market,amount,price))
			log.Println("=========================================================")
		case "sl":
			ordertype=SELL_LIMIT
			log.Println("=======输入 卖出数量=========")
			amount = input()
			log.Println("=======输入 卖出单价=========")
			price = input()
			log.Println(fmt.Sprintf("%s     限价卖出     数量:%s 单价:%s ",market,amount,price))
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
		postOrder(market, ordertype, a, p)
	}

}
func input() string {

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}
