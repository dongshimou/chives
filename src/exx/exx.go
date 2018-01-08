package exx

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

type testJson struct {
	data [7]string
}

func remove(s string) string {
	s = s[1:]
	s = s[:len(s)-1]
	return s
}
func parseSingle(raw string) *testJson {
	res := testJson{}
	raw = remove(raw)
	list := strings.Split(raw, ",")
	for i, v := range list {
		v = remove(v)
		res.data[i] = v
	}
	return &res
}
func parseMulti(raw string) []string {
	raw = remove(raw)
	return strings.Split(raw, ",")
}
func dataShow(data *testJson) {
	log.Println(fmt.Sprintf("时间:%s ==== 行情类型:%s : 价格:%s 成交量:%s",
		time.Now().Format("2006-01-02 15:04:05"), data.data[3], data.data[5], data.data[6]))
}

func Run() {

	c, _, err := websocket.DefaultDialer.Dial("wss://ws.exx/websocket", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	msg := map[string]interface{}{
		"dataType": "1_TRADE_ETH_USDT",
		"dataSize": 10,
		"action":   "ADD",
	}

	err = c.WriteJSON(msg)
	if err != nil {
		log.Println(err.Error())
	}

	go func() {
		defer c.Close()
		//first 50
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("err is ", err.Error())
		}
		log.Println(string(message))

		//list := parseMulti(string(message))
		//
		//for _, v := range list {
		//	data := parseSingle(v)
		//	dataShow(data)
		//}
		log.Println("============以上为最近的历史数据============")
		for {
			//v := testJson{}
			//err := c.ReadJSON(&v)
			//if err != nil {
			//	log.Println(err.Error())
			//}
			//log.Println(v)

			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("err is ", err.Error())
			}
			data := parseSingle(string(message))
			dataShow(data)
		}
	}()
}
