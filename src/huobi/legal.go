package huobi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	CNY_BTC  = 1
	CNY_USDT = 2

	CNY_BTC_STR  = "  CNY 与 BTC  的结算价格  "
	CNY_USDT_STR = "  CNY 与 USDT 的结算价格  "
)
const (
	LEGAL_TRADETYPE   = "0"  // 未知参数默认0
	LEGAL_CURRENTPAGE = "1"  //查第一页
	LEGAL_PAYWAY      = "2"  // 1银行卡 2支付宝 3微信
	LEGAL_COUNTRY     = "86" //中国
	LEGAL_MERCHANT    = "0"  //商人默认0
	LEGAL_ONLINE      = "1"  //0离线 1在线
	LEGAL_RANGE       = "0"  //未知参数默认0
)

func GetLegal(typeid int) {

	arg := fmt.Sprintf("https://api-otc.huobi.pro/v1/otc/trade/list/public"+
		"?coinId=%d&tradeType=%s&currentPage=%s&payWay=%s&country=%s&merchant=%s&online=%s&range=%s",
		typeid, LEGAL_TRADETYPE, LEGAL_CURRENTPAGE, LEGAL_PAYWAY, LEGAL_COUNTRY, LEGAL_MERCHANT, LEGAL_ONLINE, LEGAL_RANGE)

	res, err := http.Get(arg)
	defer res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("返回失败")
		return
	}
	raw, err := ioutil.ReadAll(res.Body)
	//log.Println(string(raw))
	data := RESLegal{}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Println("json parse error")
		return
	}
	log.Println("========================================================")
	if typeid == CNY_BTC {
		log.Println(CNY_BTC_STR)
	} else {
		log.Println(CNY_USDT_STR)
	}
	for _, v := range data.Data {
		log.Println(fmt.Sprintf(" === 价格 : %.2f, 当前量 : %12.2f  [ 当月成交 : %8d === %s  ]",
			v.Price, v.TradeCount, v.TradeMonthTimes, v.UserName))
	}
	log.Println("========================================================")
}
