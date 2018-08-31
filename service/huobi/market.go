package huobi

import (
	"net/http"
	"net/url"

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
		"kline":         getKLine(),
		"detail_merged": getTicker(),
		"depth":         getMarketDepth(),
		"trade":         getTradeDetail(),
		"detail":        getMarketDetail(),
		// "symbols":       getSymbols(),
		"currencys": getCurrencys(),
		"timestamp": getTimestamp(),
	}
)

func genericFunc(url string, query url.Values) HandleFunc {
	return func() (string, error) {
		var path = url
		var params = ""

		header := make(http.Header)
		header.Set("User-Agent", UA)

		path = url
		if query != nil {
			params = client.Map2UrlQuery(query)
			path = url + "?" + params
		}

		request := &client.HttpRequest{
			Url:    path,
			Header: header,
		}
		jsonStr, err := client.HttpGet(request)
		if err != nil {
			return "", err
		}

		return jsonStr, nil
	}
}

//------------------------------------------------------------------------------------------
// 交易API

// 获取K线数据
// 交易对, btcusdt, bccbtc......
// K线类型, 1min, 5min, 15min......
// 获取数量, [1-2000]
// return: string
func getKLine() HandleFunc {
	query := url.Values{
		"symbol": []string{"btcusdt"},
		"period": []string{"1min"},
		"size":   []string{"3"},
	}

	return genericFunc(marketAPI["kline"], query)
}

// 获取聚合行情
// 交易对, btcusdt, bccbtc......
// return: string
func getTicker() HandleFunc {
	query := url.Values{"symbol": []string{"btcusdt"}}
	return genericFunc(marketAPI["detail_merged"], query)
}

// 获取交易深度信息
// 交易对, btcusdt, bccbtc......
// Depth类型, step0、step1......stpe5 (合并深度0-5, 0时不合并)
// return: string
func getMarketDepth() HandleFunc {
	query := url.Values{
		"symbol": []string{"btcusdt"},
		"type":   []string{"step1"},
	}
	return genericFunc(marketAPI["depth"], query)
}

// 获取交易细节信息
// 交易对, btcusdt, bccbtc......
// return: string
func getTradeDetail() HandleFunc {
	query := url.Values{"symbol": []string{"btcusdt"}}
	return genericFunc(marketAPI["trade"], query)
}

// 获取Market Detail 24小时成交量数据
// 交易对, btcusdt, bccbtc......
// return: string
func getMarketDetail() HandleFunc {
	query := url.Values{"symbol": []string{"btcusdt"}}
	return genericFunc(marketAPI["detail"], query)
}

//------------------------------------------------------------------------------------------
// 公共API

// 查询系统支持的所有交易及精度
// return: string
func getSymbols() HandleFunc {
	return genericFunc(tradeAPI["symbols"], nil)
}

// 查询系统支持的所有币种
// return: string
func getCurrencys() HandleFunc {
	return genericFunc(tradeAPI["currencys"], nil)
}

// 查询系统当前时间戳
// return: string
func getTimestamp() HandleFunc {
	return genericFunc(tradeAPI["timestamp"], nil)
}
