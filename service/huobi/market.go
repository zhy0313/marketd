package huobi

import (
	"github.com/gnuos/marketd/markets"
	"github.com/gnuos/marketd/utils/client"
)

const (
	HUOBI_MARKET = "https://api.huobi.pro/market"
	HUOBI_TRADE  = "https://api.huobi.pro/v1"

	UA = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"
)

var (
	marketAPI = map[string]string{
		"kline":         HUOBI_MARKET + "/history/kline",
		"detail_merged": HUOBI_MARKET + "/detail/merged",
		"detail":        HUOBI_MARKET + "/detail",
		"tickers":       HUOBI_MARKET + "/tickers",
		"depth":         HUOBI_MARKET + "/depth",
		"trade":         HUOBI_MARKET + "/trade",
	}

	tradeAPI = map[string]string{
		"symbols":   HUOBI_TRADE + "/common/symbols",
		"currencys": HUOBI_TRADE + "/common/currencys",
		"timestamp": HUOBI_TRADE + "/common/timestamp",
	}

	handlers = map[string]HandleFunc{
		"kline":         getKLine,
		"detail_merged": getTicker,
		"depth":         getMarketDepth,
		"trade":         getTradeDetail,
		"detail":        getMarketDetail,
		"symbols":       getSymbols,
		"currencys":     getCurrencys,
		"timestamp":     getTimestamp,
	}
)

func genericFunc(url string, query map[string]string, data markets.Metric) HandleFunc {
	return func() (string, markets.Metric, error) {
		request := client.HttpRequest{
			Url:       url,
			Payload:   query,
			UserAgent: UA,
		}
		jsonStr, err := client.HttpGet(request, data)
		if err != nil {
			return "", nil, err
		}

		return jsonStr, data, nil
	}
}

//------------------------------------------------------------------------------------------
// 交易API

// 获取K线数据
// 交易对, btcusdt, bccbtc......
// K线类型, 1min, 5min, 15min......
// 获取数量, [1-2000]
// return: string
func getKLine() (string, markets.Metric, error) {
	getJson := genericFunc(
		tradeAPI["kline"],
		map[string]string{
			"symbol": "btcusdt",
			"period": "1min",
			"size":   "3",
		},
		new(KLineReturn))
	return getJson()
}

// 获取聚合行情
// 交易对, btcusdt, bccbtc......
// return: string
func getTicker() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["detail_merged"], map[string]string{"symbol": "btcusdt"}, new(TickerReturn))
	return getJson()
}

// 获取交易深度信息
// 交易对, btcusdt, bccbtc......
// Depth类型, step0、step1......stpe5 (合并深度0-5, 0时不合并)
// return: string
func getMarketDepth() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["depth"], map[string]string{"symbol": "btcusdt", "type": "step1"}, new(MarketDepthReturn))
	return getJson()
}

// 获取交易细节信息
// 交易对, btcusdt, bccbtc......
// return: string
func getTradeDetail() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["trade"], map[string]string{"symbol": "btcusdt"}, new(TradeDetailReturn))
	return getJson()
}

// 获取Market Detail 24小时成交量数据
// 交易对, btcusdt, bccbtc......
// return: string
func getMarketDetail() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["detail"], map[string]string{"symbol": "btcusdt"}, new(MarketDetailReturn))
	return getJson()
}

//------------------------------------------------------------------------------------------
// 公共API

// 查询系统支持的所有交易及精度
// return: string
func getSymbols() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["symbols"], nil, new(SymbolsReturn))
	return getJson()
}

// 查询系统支持的所有币种
// return: string
func getCurrencys() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["currencys"], nil, new(CurrencysReturn))
	return getJson()
}

// 查询系统当前时间戳
// return: string
func getTimestamp() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["timestamp"], nil, new(TimestampReturn))
	return getJson()
}
