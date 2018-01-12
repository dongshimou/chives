package huobi

import (
	"encoding/json"
	"fmt"
	"github.com/smallnest/goreq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	SPOT_Balance = MyBalances{
		map[string]float64{},
		map[string]float64{},
	}
	OTC_Balance = MyBalances{
		map[string]float64{},
		map[string]float64{},
	}
)

type MyBalances struct {
	TradeBalance  map[string]float64
	FrozenBalance map[string]float64
}

type HuobiAccount struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`    // spot , otc
	SubType string `json:"subtype"` //
	State   string `json:"state"`   // working
}

var (
	SPOT_Account int64 = 0 //现货ID
	OTC_Account  int64 = 0 //c2cID
)

func getAccountID() {

	path := "/v1/account/accounts"

	url, err := createUrl(GET, path, nil)

	if err != nil {
		log.Println("构造 url 错误")
		return
	}

	res, err := http.Get(url)

	if err != nil {
		log.Println("HTTP 错误")
		return
	}

	raw, err := ioutil.ReadAll(res.Body)
	//log.Println(string(raw))

	type tempS struct {
		Status string         `json:"status"`
		Data   []HuobiAccount `json:"data"`
	}
	temp := tempS{}
	if in := strings.Index(string(raw), "ok"); in > 0 {
		err := json.Unmarshal(raw, &temp)
		if err != nil {
			log.Println("json parse 错误")
			return
		}

		for _, v := range temp.Data {
			if v.Type == "spot" {
				SPOT_Account = v.ID
			} else if v.Type == "otc" {
				OTC_Account = v.ID
			}
		}
	}
	if SPOT_Account == 0 {
		log.Println("查询 现货 站点ID 错误")
	}
	if OTC_Account == 0 {
		log.Println("查询 结算 站点ID 错误")
	}
}

type HuobiBalance struct {
	Currency string `json:"currency"` // 类型 如 usdt
	Type     string `json:"type"`     //trade 可交易 ; frozen 冻结
	Balance  string `json:"balance"`  //余额
}

func getBalance(id int64) {
	//GET /v1/account/accounts/{account-id}/balance 查询指定账户的余额

	path := fmt.Sprintf("/v1/account/accounts/%d/balance", id)

	url, err := createUrl(GET, path, nil)

	if err != nil {
		log.Println("构造 url 错误")
		return
	}

	res, err := http.Get(url)

	if err != nil {
		log.Println("HTTP 错误")
		return
	}

	raw, err := ioutil.ReadAll(res.Body)
	//log.Println(string(raw))

	type tempData struct {
		ID    int64          `json:"id"`
		Type  string         `json:"type"`
		State string         `json:"state"`
		List  []HuobiBalance `json:"list"`
	}
	type tempRES struct {
		Status string   `json:"status"`
		Data   tempData `json:"data"`
	}

	temp := tempRES{}
	if in := strings.Index(string(raw), "ok"); in > 0 {
		err := json.Unmarshal(raw, &temp)
		if err != nil {
			log.Println("json parse 错误")
			return
		}

		processFunc := func(b *MyBalances) {
			for _, v := range temp.Data.List {
				money, _ := strconv.ParseFloat(v.Balance, 64)
				if v.Type == "trade" {
					b.TradeBalance[v.Currency] = money
				} else if v.Type == "frozen" {
					b.FrozenBalance[v.Currency] = money
				}
			}
		}

		if temp.Data.Type == "spot" {
			processFunc(&SPOT_Balance)
		} else if temp.Data.Type == "otc" {
			processFunc(&OTC_Balance)
		}
	}
}

const (
	BUY_MARKET = "buy-market"  //市价买
	SEL_MARKET = "sell-market" //市价卖
	BUY_LIMIT  = "buy-limit"   //限价买
	SELL_LIMIT = "sell-limit"  //限价卖
)
type orderResponse struct {
	Status string `json:"status"`
	Data string `json:"data"`
}
func postOrder(market, orderType string, amount, price float64) bool {

	para := map[string]interface{}{
		//api，如果使用借贷资产交易，请填写‘margin-api’
		"source": "",
		//buy-market：市价买, sell-market：市价卖, buy-limit：限价买, sell-limit：限价卖
		"type": orderType,
		//币币交易使用‘spot’账户的accountid；借贷资产交易，请使用‘margin’账户的accountid
		"account-id": fmt.Sprintf("%d", SPOT_Account),
		"symbol":     market, //交易对 btcusdt,
	}

	//amount 限价单表示下单数量，市价买单时表示买多少钱，市价卖单时表示卖多少币
	switch orderType {
	case BUY_MARKET:
		para["amount"] = fmt.Sprintf("%.4f", amount*price)
	case SEL_MARKET:
		para["amount"] = fmt.Sprintf("%.4f", amount)
	case BUY_LIMIT, SELL_LIMIT:
		para["amount"] = fmt.Sprintf("%.4f", amount)
		para["price"] = fmt.Sprintf("%.4f", price)
	}

	path := "/v1/order/orders/place"

	url, err := createUrl(POST, path, para)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_, body, errs := goreq.New().Post(url).SendStruct(para).End()
	if len(errs) != 0 {
		log.Println("HTTP errors")
		return false
	}

	log.Println(body)

	return true
}

func closeOrder(id interface{})bool{

	//POST /v1/order/orders/{order-id}/submitcancel

	para:=map[string]interface{}{
		"order-id":id,
	}

	path:=fmt.Sprintf("/v1/order/orders/%v/submitcancel",id)

	url, err := createUrl(POST, path, para)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	resJson:=orderResponse{}
	_, body, errs := goreq.New().Post(url).SendStruct(para).BindBody(&resJson).End()
	if len(errs) != 0 {
		log.Println("HTTP errors")
		return false
	}

	log.Println(body)

	return true
}