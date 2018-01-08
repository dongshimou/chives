


[huobi websocket文档](https://github.com/huobiapi/API_Docs/wiki/WS_request)
| topic 类型    | topic 语法                   | 描述                                                                                            |
| ------------- | ---------------------------- | ----------------------------------------------------------------------------------------------- |
| KLine         | market.$symbol.kline.$period | $period 可选值：{ 1min, 5min, 15min, 30min, 60min, 1day, 1mon, 1week, 1year }                   |
| Market Depth  | market.$symbol.depth.$type   | $type 可选值：{ step0, step1, step2, step3, step4, step5 } （合并深度0-5）；step0时，不合并深度 |
| Trade Detail  | market.$symbol.trade.detail  |
| Market Detail | market.$symbol.detail        |