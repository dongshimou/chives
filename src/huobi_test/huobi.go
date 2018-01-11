package huobi_test

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func Run() {

	getAccountID()

	for {
		log.Println("\n\n\n\n\n==============================")
		log.Println("====输入交易对==== 如 : ethusdt")
		market := input()
		log.Println("===输入交易方式=== 如 : =======")
		log.Println("    buy-market   //市价买")
		log.Println("    sell-market  //市价卖")
		log.Println("    buy-limit    //限价买")
		log.Println("    sell-limit   //限价卖")
		log.Println("==============================")
		ordertype := input()
		log.Println("=======输入购买量 amount=======")
		log.Println("限价单表示下单数量，")
		log.Println("市价买单时表示买多少钱，")
		log.Println("市价卖单时表示卖多少币  ")
		amount := input()
		log.Println("=======输入交易价  price=======")
		log.Println("下单价格，市价单时输入0 ")
		price := input()

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
