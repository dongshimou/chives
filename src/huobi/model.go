package huobi

import (
	"fmt"
)

type Trade struct {
	CH   string    `json:"ch"`
	TS   int64     `json:"ts"`
	Tick TradeTick `json:"tick"`
}
type TradeTick struct {
	ID   int64           `json:"id"`
	TS   int64           `json:"ts"`
	Data []TradeTickData `json:"data"`
}
type TradeTickData struct {
	ID        int64   `json:"id"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
	TS        int64   `json:"ts"`
}

const ( // TradeListData direction
	BUY  = "buy"
	SELL = "sell"
)

type TradeHello struct {
	ID      string `json:"id"`
	Subbed  string `json:"subbed"`
	TS      int64  `json:"ts"`
	Status  string `json:"status"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type RESLegal struct {
	ResponseTemplate
	Data []LegalDataInfo `json:"data"`
}
type LegalDataInfo struct {
	ID      int64  `json:"id"`
	TradeNo string `json:"tradeNo"`
	//注册国家 86 中国
	Country   int `json:"country"`
	coinID    int `json:"coinId"`
	TradeType int `json:"tradeType"`
	Merchant  int `json:"merchant"`
	//当前国家
	CurrenCY int `json:"currency"`
	//1银行卡 2支付宝 3微信
	PayMethod string `json:"payMethod"`
	UserID    int64  `json:"userId"`
	UserName  string `json:"userName"`
	//是否固定价格
	IsFixed bool `json:"isFixed"`
	//最小交易量
	MinTradeLimit float32 `json:"minTradeLimit"`
	//最大交易量
	MaxTradeLimit float32 `json:"maxTradeLimit"`
	//固定价格
	FixedPrice float32 `json:"fixedPrice"`
	//计算比例
	CalcRate float32 `json:"calcRate"`
	//交易价格
	Price float32 `json:"price"`
	//交易量
	TradeCount float32 `json:"tradeCount"`
	//是否在线
	IsOnline bool `json:"isOnline"`
	//当月交易
	TradeMonthTimes int `json:"tradeMonthTimes"`
	//当月申诉
	AppealMonthTimes int `json:"appealMonthTimes"`
	//当月胜诉
	AppealMonthWinTimes int `json:"appealMonthWinTimes"`
}
type ResponseTemplate struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	TotalCount int    `json:"totalCount"`
	PageSize   int    `json:"pageSize"`
	TotalPage  int    `json:"totalPage"`
	Success    bool   `json:"success"`
	CurrPage   int    `json:"currPage"`
}

func NewError(code int, args ...interface{}) *InnerError {
	s := InnerError{}
	return s.New(code, args...)
}

type InnerError struct {
	Code int
	Msg  string
}

func (e *InnerError) Error() string {
	return fmt.Sprintf("Code : %d , Msg : %s ", e.Code, e.Msg)
}
func (*InnerError) New(code int, args ...interface{}) *InnerError {
	return &InnerError{
		Code: code,
		Msg:  fmt.Sprintf("%v", args...),
	}
}
