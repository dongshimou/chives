package huobi

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func sendKline(c *websocket.Conn, market string) {

	//orm:=GetDBByModel(&Kline{})
	//info:=Kline{}
	//orm.Last(&info)
	//if info.Timestamp==0{
	//	info.Timestamp=1325347200
	//}
	req := map[string]interface{}{
		"req": fmt.Sprintf("market.%s.kline.1min", market),
		"id":  "id2333333",
		//"from": fmt.Sprintf("%d",info.Timestamp),
		//"from":"1515597000",
		//"to":"1515599160",
	}
	c.WriteJSON(&req)
}
