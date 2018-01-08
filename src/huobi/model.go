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
	Amount    float32 `json:"amount"`
	Price     float32 `json:"price"`
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
