package huobi

import (
	"github.com/gnuos/marketd/markets"
	"github.com/gnuos/marketd/utils"
)

const (
	HUOBI_MARKET = "https://api.huobi.pro/market"
	HUOBI_TRADE  = "https://api.huobi.pro/v1"

	UA = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"
)

type HandleFunc func(string, map[string]string) func() (string, markets.Metric)

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

	handlers = map[string]func() (string, markets.Metric, error){
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

func genericFunc(url string, query map[string]string, data markets.Metric) func() (string, markets.Metric, error) {
	return func() (string, markets.Metric, error) {
		request := utils.HttpRequest{
			Url:       url,
			Payload:   query,
			UserAgent: UA,
		}
		jsonStr, err := utils.HttpGet(request, data)
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
	data := new(KLineReturn)

	request := utils.HttpRequest{
		Url: marketAPI["kline"],
		Payload: map[string]string{
			"symbol": "btcusdt",
			"period": "1min",
			"size":   "",
		},
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

// 获取聚合行情
// 交易对, btcusdt, bccbtc......
// return: string
func getTicker() (string, markets.Metric, error) {
	data := new(TickerReturn)

	request := utils.HttpRequest{
		Url: marketAPI["detail_merged"],
		Payload: map[string]string{
			"symbol": "btcusdt",
		},
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

// 获取交易深度信息
// 交易对, btcusdt, bccbtc......
// Depth类型, step0、step1......stpe5 (合并深度0-5, 0时不合并)
// return: string
func getMarketDepth() (string, markets.Metric, error) {
	data := new(MarketDepthReturn)

	request := utils.HttpRequest{
		Url: marketAPI["depth"],
		Payload: map[string]string{
			"symbol": "btcusdt",
			"type":   "step1",
		},
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

// 获取交易细节信息
// 交易对, btcusdt, bccbtc......
// return: string
func getTradeDetail() (string, markets.Metric, error) {
	data := new(TradeDetailReturn)

	request := utils.HttpRequest{
		Url: marketAPI["trade"],
		Payload: map[string]string{
			"symbol": "btcusdt",
		},
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

// 获取Market Detail 24小时成交量数据
// 交易对, btcusdt, bccbtc......
// return: string
func getMarketDetail() (string, markets.Metric, error) {
	data := new(MarketDetailReturn)

	request := utils.HttpRequest{
		Url: marketAPI["detail"],
		Payload: map[string]string{
			"symbol": "btcusdt",
		},
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

//------------------------------------------------------------------------------------------
// 公共API

// 查询系统支持的所有交易及精度
// return: string
func getSymbols() (string, markets.Metric, error) {
	data := new(SymbolsReturn)

	request := utils.HttpRequest{
		Url:       tradeAPI["symbols"],
		Payload:   nil,
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

// 查询系统支持的所有币种
// return: string
func getCurrencys() (string, markets.Metric, error) {
	data := new(CurrencysReturn)

	request := utils.HttpRequest{
		Url:       tradeAPI["currencys"],
		Payload:   nil,
		UserAgent: UA,
	}

	jsonStr, err := utils.HttpGet(request, data)
	if err != nil {
		return "", nil, err
	}

	return jsonStr, data, nil
}

// 查询系统当前时间戳
// return: string
func getTimestamp() (string, markets.Metric, error) {
	getJson := genericFunc(tradeAPI["timestamp"], nil, new(TimestampReturn))
	return getJson()
}
