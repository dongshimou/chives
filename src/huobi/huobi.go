package huobi

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"math"
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

func Run() {
	getAccountID()
	postOrder("ethusdt", SEL_MARKET, 0.002, 0)
	market := ""
	for {
		c, err := initWS()
		market = input()
		//market = "ethusdt"
		addMarket(market)
		InitDB(market)
		log.Println("所查行情为 :", market)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = startWS(c, market)
		if err == nil {
			break
		} else {
			log.Println("==========启动失败正在重新启动==========\n\n\n")
		}
	}
	retryTrade := func() {
		c, err := initWS()
		InitDB(market)
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = startWS(c, market)
		if err != nil {
			log.Println("==========启动失败正在重新启动==========\n\n\n")
		}
	}
	//检查是否关闭
	go func() {
		for {
			if isDetailClose() {
				log.Println("即时行情异常关闭,正在重启")
				go retryTrade()
			}
			time.Sleep(time.Second * 10)
		}
	}()
	priceList := map[int]float64{}
	minuteCount := 0
	for {
		setCurrentInit()
		time.Sleep(time.Minute)
		getLegal(CNY_USDT)
		totalmax, totalmin := getTotalMaxMinPrice()
		currmax, currmin := getCurrentMaxMinPrice()
		currAmount := getCurrentAmount()
		log.Println(fmt.Sprintf("========================当前 %s 的价格======================", market))
		log.Println(fmt.Sprintf("===历史的最高价 : %14.4f ===========历史的最低价 : %14.4f ===", totalmax, totalmin))
		log.Println(fmt.Sprintf("===上分钟最高价 : %14.4f ===========上分钟最低价 : %14.4f ===", currmax, currmin))
		log.Println(fmt.Sprintf("=======================上分钟成交量: %32.4f ===============", currAmount))

		temp := isUpOrDown()
		up, bala, down := getUpDownCount()
		log.Println(fmt.Sprintf("上涨 %14.0f == 相等 %14.0f == 下跌 %14.0f ", up, bala, down))
		if temp > 0 {
			log.Println(" ===== 上涨多 =====")
		} else if temp == 0 {
			log.Println(" ===== 震荡中 =====")
		} else {
			log.Println(" ===== 下跌多 =====")
		}
		priceList[minuteCount] = getAveragePrice()
		log.Println("============================当前均价====", priceList[minuteCount])

		minuteCount++
		check := func(n, o float64) int {
			rage := float64(0.0003)
			temp := n - o
			if math.Abs(temp) < rage {
				return 0
			} else if temp < 0 {
				return -1
			} else {
				return 1
			}
		}
		if minuteCount > 5 {
			//(priceList[minuteCount-1] - priceList[minuteCount-2])/priceList[minuteCount-2]>0.0001
			c1 := check(priceList[minuteCount-1], priceList[minuteCount-2])
			c2 := check(priceList[minuteCount-2], priceList[minuteCount-3])
			c3 := check(priceList[minuteCount-3], priceList[minuteCount-4])
			if c1 == 1 && c2 == 1 && c3 == 1 {
				log.Println("================上涨中================")
				log.Println("================上涨中================")
				log.Println("================上涨中================")
			} else if c1 == -1 && c2 == -1 && c3 == -1 {
				log.Println("================下跌中================")
				log.Println("================下跌中================")
				log.Println("================下跌中================")
			}

			p1 := getLastPrice()
			p2 := priceList[minuteCount-1]
			log.Println(fmt.Sprintf("当前价格: %14.4f   上分钟均价: %14.4f", p1, p2))

			if p1 > p2 {
				log.Println(fmt.Sprintf(" 上涨幅度 %%%.2f", (p1-p2)/p2*100))
			} else {
				log.Println(fmt.Sprintf(" 下跌幅度 %%%.2f", (p2-p1)/p2*100))
			}
		}
	}
}
