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

func getSPOTBalance(symbol string)float64{
	return SPOT_Balance.TradeBalance[symbol]
}
func getOTCBalance(symbol string)float64{
	return OTC_Balance.TradeBalance[symbol]
}
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

//下单
func createOrder(market, orderType string, amount, price float64) bool {

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

type OrderResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type orderInfoResponse struct {
	ID              int64  `json:"id"`
	Symbol          string `json:"symbol"`
	AccountID       int64  `json:"account-id"`
	Amount          string `json:"amount"` //float64
	Price           string `json:"price"`  //float64
	CreatedAt       int64  `json:"created-at"`
	Type            string `json:"type"`
	FieldAmount     string `json:"field-amount"`      //float64
	FieldCashAmount string `json:"field-cash-amount"` //float64
	FieldFees       string `json:"field-fees"`        //float64
	FinishedAt      int64  `json:"finished-at"`
	UserID          int64  `json:"user-id"`
	Source          string `json:"source"`
	State           string `json:"state"`
	CanceledAt      int64  `json:"canceled-at"`
	Exchange        string `json:"exchange"`
	Batch           string `json:"batch"`
}
//GET /v1/order/orders/{order-id} 查询某个订单详情
func getOrderInfo(id interface{}) {

	para := map[string]interface{}{
		"order-id": id,
	}
	path := fmt.Sprintf("/v1/order/orders/%v", id)

	url, err := createUrl(POST, path, para)
	if err != nil {

	}
	resJson := OrderResponse{}
	info := orderInfoResponse{}
	resJson.Data = info
	_, _, errs := goreq.New().Get(url).BindBody(&resJson).End()
	if len(errs) != 0 {
		log.Println("HTTP errors")
		return
	}
	status := "订单状态:"
	switch info.State {
	// pre-submitted 准备提交,
	case "pre-submitted":
		status += " 准备提交 "
		// submitting , submitted 已提交,
	case "submitting", "submitted":
		status += " 已提交 "
		// partial-filled 部分成交,
	case "partial-filled":
		status += " 部分成交 "
		// partial-canceled 部分成交撤销,
	case "partial-canceled":
		status += " 部分成交撤销 "
	// filled 完全成交,
	case "filled":
		status += " 完全成交 "
	// canceled 已撤销
	case "canceled":
		status += " 已撤销 "
	}
	log.Println(status + fmt.Sprintf(" 交易对:%s 已成交数量:%s 已成交金额:%s 订单ID:%d 订单价格:%s ",
		info.Type, info.FieldAmount, info.FieldCashAmount, info.ID, info.Price))
}

//GET /v1/order/orders 查询当前委托、历史委托
func queryHistoryOrder(market string) {
	/*
	GET /v1/order/orders 查询当前委托、历史委托
	参数名称	是否必须	类型	描述	默认值	取值范围
	   symbol	true	string	交易对		btcusdt, bccbtc, rcneth ...
	   types	false	string	查询的订单类型组合，使用','分割		buy-market：市价买, sell-market：市价卖, buy-limit：限价买, sell-limit：限价卖
	   start-date	false	string	查询开始日期, 日期格式yyyy-mm-dd
	   end-date	false	string	查询结束日期, 日期格式yyyy-mm-dd
	   states	true	string	查询的订单状态组合，使用','分割		pre-submitted 准备提交, submitted 已提交, partial-filled 部分成交, partial-canceled 部分成交撤销, filled 完全成交, canceled 已撤销
	   from	false	string	查询起始 ID
	   direct	false	string	查询方向		prev 向前，next 向后
	   size	false	string	查询记录大小
	*/

	para:=map[string]interface{}{
		"symbol":market,
		"states":"submitted",
	}
	path:="/v1/order/orders"

	url,err:=createUrl(GET,path,para)
	if err!=nil{
		log.Println(err.Error())
		return
	}

	resJson := OrderResponse{}
	resJson.Data=orderInfoResponse{}
	_, _, errs := goreq.New().Get(url).BindBody(&resJson).End()
	if len(errs) != 0 {
		log.Println("HTTP errors")
		return
	}
}

// 取消订单
func closeOrder(id interface{}) bool {

	//POST /v1/order/orders/{order-id}/submitcancel

	para := map[string]interface{}{
		"order-id": id,
	}

	path := fmt.Sprintf("/v1/order/orders/%v/submitcancel", id)

	url, err := createUrl(POST, path, para)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	resJson := OrderResponse{}
	resJson.Data=""
	_, body, errs := goreq.New().Post(url).SendStruct(para).BindBody(&resJson).End()
	if len(errs) != 0 {
		log.Println("HTTP errors")
		return false
	}

	log.Println(body)

	return true
}
